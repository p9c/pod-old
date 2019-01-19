package node

import (
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime/pprof"

	"git.parallelcoin.io/pod/lib/blockchain/indexers"
	"git.parallelcoin.io/pod/lib/database"
)

const (
	// blockDbNamePrefix is the prefix for the block database name.  The database type is appended to this value to form the full block database name.
	blockDbNamePrefix = "blocks"
)

var (
	cfg *Config
)

// winServiceMain is only invoked on Windows.  It detects when pod is running as a service and reacts accordingly.
var winServiceMain func() (bool, error)

// Main is the real main function for pod.  It is necessary to work around the fact that deferred functions do not run when os.Exit() is called.  The optional serverChan parameter is mainly used by the service code to be notified with the server once it is setup so it can gracefully stop it when requested from the service control manager.
func Main(c *Config, serverChan chan<- *server) (err error) {
	cfg = c
	// Load configuration and parse command line.  This function also initializes logging and configures it accordingly.
	// tcfg, _, err := loadConfig()
	// if err != nil {
	// 	return err
	// }
	// cfg = tcfg
	// defer func() {
	// 	if logRotator != nil {
	// 		logRotator.Close()
	// 	}
	// }()
	// Get a channel that will be closed when a shutdown signal has been triggered either from an OS signal such as SIGINT (Ctrl+C) or from another subsystem such as the RPC server.
	interrupt := interruptListener()
	defer Log.Info.Print("shutdown complete")
	// Show version at startup.
	Log.Infof.Print("version %s", Version())
	// Enable http profiling server if requested.
	if cfg.Profile != "" {
		go func() {
			listenAddr := net.JoinHostPort("", cfg.Profile)
			Log.Infof.Print("profile server listening on %s", listenAddr)
			profileRedirect := http.RedirectHandler("/debug/pprof",
				http.StatusSeeOther)
			http.Handle("/", profileRedirect)
			Log.Errorf.Print("%v", http.ListenAndServe(listenAddr, nil))
		}()
	}
	// Write cpu profile if requested.
	if cfg.CPUProfile != "" {
		var f *os.File
		f, err = os.Create(cfg.CPUProfile)
		if err != nil {
			Log.Errorf.Print("unable to create cpu profile: %v", err)
			return
		}
		pprof.StartCPUProfile(f)
		defer f.Close()
		defer pprof.StopCPUProfile()
	}
	// Perform upgrades to pod as new versions require it.
	if err = doUpgrades(); err != nil {
		Log.Errorf.Print("%v", err)
		return
	}
	// Return now if an interrupt signal was triggered.
	if interruptRequested(interrupt) {
		return nil
	}
	// Load the block database.
	var db database.DB
	db, err = loadBlockDB()
	if err != nil {
		Log.Errorf.Print("%v", err)
		return
	}
	defer func() {
		// Ensure the database is sync'd and closed on shutdown.
		Log.Info <- "gracefully shutting down the database..."
		db.Close()
	}()
	// Return now if an interrupt signal was triggered.
	if interruptRequested(interrupt) {
		return nil
	}
	// Drop indexes and exit if requested. NOTE: The order is important here because dropping the tx index also drops the address index since it relies on it.
	if cfg.DropAddrIndex {
		if err = indexers.DropAddrIndex(db, interrupt); err != nil {
			Log.Errorf.Print("%v", err)
			return
		}
		return nil
	}
	if cfg.DropTxIndex {
		if err = indexers.DropTxIndex(db, interrupt); err != nil {
			Log.Errorf.Print("%v", err)
			return
		}
		return nil
	}
	if cfg.DropCfIndex {
		if err := indexers.DropCfIndex(db, interrupt); err != nil {
			Log.Errorf.Print("%v", err)
			return err
		}
		return nil
	}
	// Create server and start it.
	server, err := newServer(cfg.Listeners, db, ActiveNetParams.Params, interrupt, cfg.Algo)
	if err != nil {
		// TODO: this logging could do with some beautifying.
		Log.Errorf.Print("unable to start server on %v: %v", cfg.Listeners, err)
		return err
	}
	defer func() {
		Log.Info <- "gracefully shutting down the server..."
		server.Stop()
		server.WaitForShutdown()
		Log.Info <- "server shutdown complete"
	}()
	server.Start()
	if serverChan != nil {
		serverChan <- server
	}
	// Wait until the interrupt signal is received from an OS signal or shutdown is requested through one of the subsystems such as the RPC server.
	<-interrupt
	return nil
}

// removeRegressionDB removes the existing regression test database if running in regression test mode and it already exists.
func removeRegressionDB(dbPath string) error {
	// Don't do anything if not in regression test mode.
	if !cfg.RegressionTest {
		return nil
	}
	// Remove the old regression test database if it already exists.
	fi, err := os.Stat(dbPath)
	if err == nil {
		Log.Infof.Print("removing regression test database from '%s'", dbPath)
		if fi.IsDir() {
			err := os.RemoveAll(dbPath)
			if err != nil {
				return err
			}
		} else {
			err := os.Remove(dbPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// dbPath returns the path to the block database given a database type.
func blockDbPath(dbType string) string {
	// The database name is based on the database type.
	dbName := blockDbNamePrefix + "_" + dbType
	if dbType == "sqlite" {
		dbName = dbName + ".db"
	}
	dbPath := filepath.Join(cfg.DataDir, dbName)
	return dbPath
}

// warnMultipleDBs shows a warning if multiple block database types are detected. This is not a situation most users want.  It is handy for development however to support multiple side-by-side databases.
func warnMultipleDBs() {
	// This is intentionally not using the known db types which depend on the database types compiled into the binary since we want to detect legacy db types as well.
	dbTypes := []string{"ffldb", "leveldb", "sqlite"}
	duplicateDbPaths := make([]string, 0, len(dbTypes)-1)
	for _, dbType := range dbTypes {
		if dbType == cfg.DbType {
			continue
		}
		// Store db path as a duplicate db if it exists.
		dbPath := blockDbPath(dbType)
		if FileExists(dbPath) {
			duplicateDbPaths = append(duplicateDbPaths, dbPath)
		}
	}
	// Warn if there are extra databases.
	if len(duplicateDbPaths) > 0 {
		selectedDbPath := blockDbPath(cfg.DbType)
		Log.Warnf.Print("\nThere are multiple block chain databases using different database types.\n"+
			"You probably don't want to waste disk space by having more than one.\n"+
			"Your current database is located at [%v].\n"+
			"The additional database is located at %v", selectedDbPath, duplicateDbPaths)
	}
}

// loadBlockDB loads (or creates when needed) the block database taking into account the selected database backend and returns a handle to it.  It also additional logic such warning the user if there are multiple databases which consume space on the file system and ensuring the regression test database is clean when in regression test mode.
func loadBlockDB() (database.DB, error) {
	// The memdb backend does not have a file path associated with it, so handle it uniquely.  We also don't want to worry about the multiple database type warnings when running with the memory database.
	if cfg.DbType == "memdb" {
		Log.Info <- "creating block database in memory"
		db, err := database.Create(cfg.DbType)
		if err != nil {
			return nil, err
		}
		return db, nil
	}
	warnMultipleDBs()
	// The database name is based on the database type.
	dbPath := blockDbPath(cfg.DbType)
	// The regression test is special in that it needs a clean database for each run, so remove it now if it already exists.
	removeRegressionDB(dbPath)
	Log.Infof.Print("loading block database from '%s'", dbPath)
	db, err := database.Open(cfg.DbType, dbPath, ActiveNetParams.Net)
	if err != nil {
		// Return the error if it's not because the database doesn't exist.
		if dbErr, ok := err.(database.Error); !ok || dbErr.ErrorCode !=
			database.ErrDbDoesNotExist {
			return nil, err
		}
		// Create the db if it does not exist.
		err = os.MkdirAll(cfg.DataDir, 0700)
		if err != nil {
			return nil, err
		}
		db, err = database.Create(cfg.DbType, dbPath, ActiveNetParams.Net)
		if err != nil {
			return nil, err
		}
	}
	Log.Info <- "Block database loaded"
	return db, nil
}

// func PreMain() {
// 	// Use all processor cores.
// 	runtime.GOMAXPROCS(runtime.NumCPU())
// 	// Block and transaction processing can cause bursty allocations.  This limits the garbage collector from excessively overallocating during bursts.  This value was arrived at with the help of profiling live usage.
// 	debug.SetGCPercent(10)
// 	// Up some limits.
// 	if err := limits.SetLimits(); err != nil {
// 		fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
// 		os.Exit(1)
// 	}
// 	// Call serviceMain on Windows to handle running as a service.  When the return isService flag is true, exit now since we ran as a service.  Otherwise, just fall through to normal operation.
// 	if runtime.GOOS == "windows" {
// 		isService, err := winServiceMain()
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 		if isService {
// 			os.Exit(0)
// 		}
// 	}
// 	// Work around defer not working after os.Exit()
// 	if err := Main(nil); err != nil {
// 		os.Exit(1)
// 	}
// }