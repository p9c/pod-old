// Copyright (c) 2013-2017 The btcsuite developers

package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/parallelcointeam/mod/internal/cfgutil"
	"github.com/parallelcointeam/mod/internal/legacy/keystore"
	"github.com/parallelcointeam/mod/netparams"
	"github.com/parallelcointeam/mod/wallet"
	"github.com/parallelcointeam/pod/btcutil"
)

const (
	defaultCAFilename       = "mod.cert"
	defaultConfigFilename   = "mod.conf"
	defaultLogLevel         = "info"
	defaultLogDirname       = "logs"
	defaultLogFilename      = "mod.log"
	defaultRPCMaxClients    = 10
	defaultRPCMaxWebsockets = 25

	walletDbName = "wallet.db"
)

var (
	podDefaultCAFile   = filepath.Join(btcutil.AppDataDir("mod", false), "rpc.cert")
	defaultAppDataDir  = btcutil.AppDataDir("mod", false)
	defaultConfigFile  = filepath.Join(defaultAppDataDir, defaultConfigFilename)
	defaultRPCKeyFile  = filepath.Join(defaultAppDataDir, "rpc.key")
	defaultRPCCertFile = filepath.Join(defaultAppDataDir, "rpc.cert")
	defaultLogDir      = filepath.Join(defaultAppDataDir, defaultLogDirname)
)

type config struct {
	// General application behavior
	ConfigFile    *cfgutil.ExplicitString `short:"C" long:"configfile" description:"Path to configuration file"`
	ShowVersion   bool                    `short:"V" long:"version" description:"Display version information and exit"`
	Create        bool                    `long:"create" description:"Create the wallet if it does not exist"`
	CreateTemp    bool                    `long:"createtemp" description:"Create a temporary simulation wallet (pass=password) in the data directory indicated; must call with --datadir"`
	AppDataDir    *cfgutil.ExplicitString `short:"A" long:"appdata" description:"Application data directory for wallet config, databases and logs"`
	TestNet3      bool                    `long:"testnet" description:"Use the test Bitcoin network (version 3) (default mainnet)"`
	SimNet        bool                    `long:"simnet" description:"Use the simulation test network (default mainnet)"`
	NoInitialLoad bool                    `long:"noinitialload" description:"Defer wallet creation/opening on startup and enable loading wallets over RPC"`
	DebugLevel    string                  `short:"d" long:"debuglevel" description:"Logging level {trace, debug, info, warn, error, critical}"`
	LogDir        string                  `long:"logdir" description:"Directory to log output."`
	Profile       string                  `long:"profile" description:"Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536"`
	GUI           bool                    `long:"gui" description:"Launch GUI"`
	// Wallet options
	WalletPass string `long:"walletpass" default-mask:"-" description:"The public wallet password -- Only required if the wallet was created with one"`

	// RPC client options
	RPCConnect      string                  `short:"c" long:"rpcconnect" description:"Hostname/IP and port of pod RPC server to connect to (default localhost:11048, testnet: localhost:21048, simnet: localhost:41048)"`
	CAFile          *cfgutil.ExplicitString `long:"cafile" description:"File containing root certificates to authenticate a TLS connections with pod"`
	EnableClientTLS bool                    `long:"clienttls" description:"Enable TLS for the RPC client"`
	PodUsername     string                  `long:"podusername" description:"Username for pod authentication"`
	PodPassword     string                  `long:"podpassword" default-mask:"-" description:"Password for pod authentication"`
	Proxy           string                  `long:"proxy" description:"Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	ProxyUser       string                  `long:"proxyuser" description:"Username for proxy server"`
	ProxyPass       string                  `long:"proxypass" default-mask:"-" description:"Password for proxy server"`

	// SPV client options
	// UseSPV       bool          `long:"usespv" description:"Enables the experimental use of SPV rather than RPC for chain synchronization"`
	AddPeers     []string      `short:"a" long:"addpeer" description:"Add a peer to connect with at startup"`
	ConnectPeers []string      `long:"connect" description:"Connect only to the specified peers at startup"`
	MaxPeers     int           `long:"maxpeers" description:"Max number of inbound and outbound peers"`
	BanDuration  time.Duration `long:"banduration" description:"How long to ban misbehaving peers.  Valid time units are {s, m, h}.  Minimum 1 second"`
	BanThreshold uint32        `long:"banthreshold" description:"Maximum allowed ban score before disconnecting and banning misbehaving peers."`

	// RPC server options
	//
	// The legacy server is still enabled by default (and eventually will be
	// replaced with the experimental server) so prepare for that change by
	// renaming the struct fields (but not the configuration options).
	//
	// Usernames can also be used for the consensus RPC client, so they
	// aren't considered legacy.
	RPCCert                *cfgutil.ExplicitString `long:"rpccert" description:"File containing the certificate file"`
	RPCKey                 *cfgutil.ExplicitString `long:"rpckey" description:"File containing the certificate key"`
	OneTimeTLSKey          bool                    `long:"onetimetlskey" description:"Generate a new TLS certpair at startup, but only write the certificate to disk"`
	EnableServerTLS        bool                    `long:"servertls" description:"Enable TLS for the RPC server"`
	LegacyRPCListeners     []string                `long:"rpclisten" description:"Listen for legacy RPC connections on this interface/port (default port: 11046, testnet: 21046, simnet: 41046)"`
	LegacyRPCMaxClients    int64                   `long:"rpcmaxclients" description:"Max number of legacy RPC clients for standard connections"`
	LegacyRPCMaxWebsockets int64                   `long:"rpcmaxwebsockets" description:"Max number of legacy RPC websocket connections"`
	Username               string                  `short:"u" long:"username" description:"Username for legacy RPC and pod authentication (if podusername is unset)"`
	Password               string                  `short:"P" long:"password" default-mask:"-" description:"Password for legacy RPC and pod authentication (if podpassword is unset)"`

	// EXPERIMENTAL RPC server options
	//
	// These options will change (and require changes to config files, etc.)
	// when the new gRPC server is enabled.
	ExperimentalRPCListeners []string `long:"experimentalrpclisten" description:"Listen for RPC connections on this interface/port"`

	// Deprecated options
	DataDir *cfgutil.ExplicitString `short:"b" long:"datadir" default-mask:"-" description:"DEPRECATED -- use appdata instead"`
}

// cleanAndExpandPath expands environement variables and leading ~ in the
// passed path, cleans the result, and returns it.
func cleanAndExpandPath(path string) string {
	// NOTE: The os.ExpandEnv doesn't work with Windows cmd.exe-style
	// %VARIABLE%, but they variables can still be expanded via POSIX-style
	// $VARIABLE.
	path = os.ExpandEnv(path)

	if !strings.HasPrefix(path, "~") {
		return filepath.Clean(path)
	}

	// Expand initial ~ to the current user's home directory, or ~otheruser
	// to otheruser's home directory.  On Windows, both forward and backward
	// slashes can be used.
	path = path[1:]

	var pathSeparators string
	if runtime.GOOS == "windows" {
		pathSeparators = string(os.PathSeparator) + "/"
	} else {
		pathSeparators = string(os.PathSeparator)
	}

	userName := ""
	if i := strings.IndexAny(path, pathSeparators); i != -1 {
		userName = path[:i]
		path = path[i:]
	}

	homeDir := ""
	var u *user.User
	var err error
	if userName == "" {
		u, err = user.Current()
	} else {
		u, err = user.Lookup(userName)
	}
	if err == nil {
		homeDir = u.HomeDir
	}
	// Fallback to CWD if user lookup fails or user has no home directory.
	if homeDir == "" {
		homeDir = "."
	}

	return filepath.Join(homeDir, path)
}

// validLogLevel returns whether or not logLevel is a valid debug log level.
func validLogLevel(logLevel string) bool {
	switch logLevel {
	case "trace":
		fallthrough
	case "debug":
		fallthrough
	case "info":
		fallthrough
	case "warn":
		fallthrough
	case "error":
		fallthrough
	case "critical":
		return true
	}
	return false
}

// supportedSubsystems returns a sorted slice of the supported subsystems for
// logging purposes.
func supportedSubsystems() []string {
	// Convert the subsystemLoggers map keys to a slice.
	subsystems := make([]string, 0, len(subsystemLoggers))
	for subsysID := range subsystemLoggers {
		subsystems = append(subsystems, subsysID)
	}

	// Sort the subsytems for stable display.
	sort.Strings(subsystems)
	return subsystems
}

// parseAndSetDebugLevels attempts to parse the specified debug level and set
// the levels accordingly.  An appropriate error is returned if anything is
// invalid.
func parseAndSetDebugLevels(debugLevel string) error {
	// When the specified string doesn't have any delimters, treat it as
	// the log level for all subsystems.
	if !strings.Contains(debugLevel, ",") && !strings.Contains(debugLevel, "=") {
		// Validate debug log level.
		if !validLogLevel(debugLevel) {
			str := "The specified debug level [%v] is invalid"
			return fmt.Errorf(str, debugLevel)
		}

		// Change the logging level for all subsystems.
		setLogLevels(debugLevel)

		return nil
	}

	// Split the specified string into subsystem/level pairs while detecting
	// issues and update the log levels accordingly.
	for _, logLevelPair := range strings.Split(debugLevel, ",") {
		if !strings.Contains(logLevelPair, "=") {
			str := "The specified debug level contains an invalid " +
				"subsystem/level pair [%v]"
			return fmt.Errorf(str, logLevelPair)
		}

		// Extract the specified subsystem and log level.
		fields := strings.Split(logLevelPair, "=")
		subsysID, logLevel := fields[0], fields[1]

		// Validate subsystem.
		if _, exists := subsystemLoggers[subsysID]; !exists {
			str := "The specified subsystem [%v] is invalid -- " +
				"supported subsytems %v"
			return fmt.Errorf(str, subsysID, supportedSubsystems())
		}

		// Validate log level.
		if !validLogLevel(logLevel) {
			str := "The specified debug level [%v] is invalid"
			return fmt.Errorf(str, logLevel)
		}

		setLogLevel(subsysID, logLevel)
	}

	return nil
}

// loadConfig initializes and parses the config using a config file and command
// line options.
//
// The configuration proceeds as follows:
//      1) Start with a default config with sane settings
//      2) Pre-parse the command line to check for an alternative config file
//      3) Load configuration file overwriting defaults with any specified options
//      4) Parse CLI options and overwrite/add any specified options
//
// The above results in btcwallet functioning properly without any config
// settings while still allowing the user to override settings with config files
// and command line options.  Command line options always take precedence.
func loadConfig() (*config, []string, error) {
	// Default config.
	cfg := config{
		DebugLevel:             defaultLogLevel,
		ConfigFile:             cfgutil.NewExplicitString(defaultConfigFile),
		AppDataDir:             cfgutil.NewExplicitString(defaultAppDataDir),
		LogDir:                 defaultLogDir,
		WalletPass:             wallet.InsecurePubPassphrase,
		CAFile:                 cfgutil.NewExplicitString(""),
		RPCKey:                 cfgutil.NewExplicitString(defaultRPCKeyFile),
		RPCCert:                cfgutil.NewExplicitString(defaultRPCCertFile),
		LegacyRPCMaxClients:    defaultRPCMaxClients,
		LegacyRPCMaxWebsockets: defaultRPCMaxWebsockets,
		DataDir:                cfgutil.NewExplicitString(defaultAppDataDir),
		// UseSPV:                 false,
		AddPeers:     []string{},
		ConnectPeers: []string{},
		// MaxPeers:               sac.MaxPeers,
		// BanDuration:            sac.BanDuration,
		// BanThreshold:           sac.BanThreshold,
	}

	// Pre-parse the command line options to see if an alternative config
	// file or the version flag was specified.
	preCfg := cfg
	preParser := flags.NewParser(&preCfg, flags.Default)
	_, err := preParser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			preParser.WriteHelp(os.Stderr)
		}
		return nil, nil, err
	}

	// Show the version and exit if the version flag was specified.
	funcName := "loadConfig"
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	usageMessage := fmt.Sprintf("Use %s -h to show usage", appName)
	if preCfg.ShowVersion {
		fmt.Println(appName, "version", version())
		os.Exit(0)
	}

	// Load additional config from file.
	var configFileError error
	parser := flags.NewParser(&cfg, flags.Default)
	configFilePath := preCfg.ConfigFile.Value
	if preCfg.ConfigFile.ExplicitlySet() {
		configFilePath = cleanAndExpandPath(configFilePath)
	} else {
		appDataDir := preCfg.AppDataDir.Value
		if !preCfg.AppDataDir.ExplicitlySet() && preCfg.DataDir.ExplicitlySet() {
			appDataDir = cleanAndExpandPath(preCfg.DataDir.Value)
		}
		if appDataDir != defaultAppDataDir {
			configFilePath = filepath.Join(appDataDir, defaultConfigFilename)
		}
	}
	err = flags.NewIniParser(parser).ParseFile(configFilePath)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			fmt.Fprintln(os.Stderr, err)
			parser.WriteHelp(os.Stderr)
			return nil, nil, err
		}
		configFileError = err
	}

	// Parse command line options again to ensure they take precedence.
	remainingArgs, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return nil, nil, err
	}

	// // Check deprecated aliases.  The new options receive priority when both
	// // are changed from the default.
	// if cfg.DataDir.ExplicitlySet() {
	// 	fmt.Fprintln(os.Stderr, "datadir option has been replaced by "+
	// 		"appdata -- please update your config")
	// 	if !cfg.AppDataDir.ExplicitlySet() {
	// 		cfg.AppDataDir.Value = cfg.DataDir.Value
	// 	}
	// }

	// If an alternate data directory was specified, and paths with defaults
	// relative to the data dir are unchanged, modify each path to be
	// relative to the new data dir.
	if cfg.AppDataDir.ExplicitlySet() {
		cfg.AppDataDir.Value = cleanAndExpandPath(cfg.AppDataDir.Value)
		if !cfg.RPCKey.ExplicitlySet() {
			cfg.RPCKey.Value = filepath.Join(cfg.AppDataDir.Value, "rpc.key")
		}
		if !cfg.RPCCert.ExplicitlySet() {
			cfg.RPCCert.Value = filepath.Join(cfg.AppDataDir.Value, "rpc.cert")
		}
	}

	if _, err := os.Stat(cfg.DataDir.Value); os.IsNotExist(err) {
		// Create the destination directory if it does not exists
		err = os.MkdirAll(cfg.DataDir.Value, 0700)
		if err != nil {
			fmt.Println("ERROR", err)
			return nil, nil, err
		}
	}

	var generatedRPCPass, generatedRPCUser string

	if _, err := os.Stat(cfg.ConfigFile.Value); os.IsNotExist(err) {

		// If we can find a pod.conf in the standard location, copy
		// copy the rpcuser and rpcpassword and TLS setting
		c := cleanAndExpandPath("~/.pod/pod.conf")
		// fmt.Println("server config path:", c)
		// _, err := os.Stat(c)
		// fmt.Println(err)
		// fmt.Println(os.IsNotExist(err))
		if _, err := os.Stat(c); err == nil {
			fmt.Println("Creating config from pod config")

			createDefaultConfigFile(cfg.ConfigFile.Value, c, cleanAndExpandPath("~/.pod"),
				cfg.AppDataDir.Value)
		} else {
			var bb bytes.Buffer
			bb.Write(sampleModConf)

			fmt.Println("Writing config file:", cfg.ConfigFile.Value)
			dest, err := os.OpenFile(cfg.ConfigFile.Value,
				os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				fmt.Println("ERROR", err)
				return nil, nil, err
			}
			defer dest.Close()

			// We generate a random user and password
			randomBytes := make([]byte, 20)
			_, err = rand.Read(randomBytes)
			if err != nil {
				return nil, nil, err
			}
			generatedRPCUser = base64.StdEncoding.EncodeToString(randomBytes)

			_, err = rand.Read(randomBytes)
			if err != nil {
				return nil, nil, err
			}
			generatedRPCPass = base64.StdEncoding.EncodeToString(randomBytes)

			// We copy every line from the sample config file to the destination,
			// only replacing the two lines for rpcuser and rpcpass
			//
			var line string
			reader := bufio.NewReader(&bb)
			for err != io.EOF {
				line, err = reader.ReadString('\n')
				if err != nil && err != io.EOF {
					return nil, nil, err
				}
				if !strings.Contains(line, "podusername=") && !strings.Contains(line, "podpassword=") {

					if strings.Contains(line, "username=") {
						line = "username=" + generatedRPCUser + "\n"
					} else if strings.Contains(line, "password=") {
						line = "password=" + generatedRPCPass + "\n"
					}
				}
				_, _ = generatedRPCPass, generatedRPCUser

				if _, err := dest.WriteString(line); err != nil {
					return nil, nil, err
				}
			}
		}
	}

	// Choose the active network params based on the selected network.
	// Multiple networks can't be selected simultaneously.
	numNets := 0
	if cfg.TestNet3 {
		activeNet = &netparams.TestNet3Params
		numNets++
	}
	if cfg.SimNet {
		activeNet = &netparams.SimNetParams
		numNets++
	}
	if numNets > 1 {
		str := "%s: The testnet and simnet params can't be used " +
			"together -- choose one"
		err := fmt.Errorf(str, "loadConfig")
		fmt.Fprintln(os.Stderr, err)
		parser.WriteHelp(os.Stderr)
		return nil, nil, err
	}

	// Append the network type to the log directory so it is "namespaced"
	// per network.
	cfg.LogDir = cleanAndExpandPath(cfg.LogDir)
	cfg.LogDir = filepath.Join(cfg.LogDir, activeNet.Params.Name)

	// Special show command to list supported subsystems and exit.
	if cfg.DebugLevel == "show" {
		fmt.Println("Supported subsystems", supportedSubsystems())
		os.Exit(0)
	}

	// Initialize log rotation.  After log rotation has been initialized, the
	// logger variables may be used.
	initLogRotator(filepath.Join(cfg.LogDir, defaultLogFilename))

	// Parse, validate, and set debug log level(s).
	if err := parseAndSetDebugLevels(cfg.DebugLevel); err != nil {
		err := fmt.Errorf("%s: %v", "loadConfig", err.Error())
		fmt.Fprintln(os.Stderr, err)
		parser.WriteHelp(os.Stderr)
		return nil, nil, err
	}

	// Exit if you try to use a simulation wallet with a standard
	// data directory.
	if !(cfg.AppDataDir.ExplicitlySet() || cfg.DataDir.ExplicitlySet()) && cfg.CreateTemp {
		fmt.Fprintln(os.Stderr, "Tried to create a temporary simulation "+
			"wallet, but failed to specify data directory!")
		os.Exit(0)
	}

	// Exit if you try to use a simulation wallet on anything other than
	// simnet or testnet3.
	if !cfg.SimNet && cfg.CreateTemp {
		fmt.Fprintln(os.Stderr, "Tried to create a temporary simulation "+
			"wallet for network other than simnet!")
		os.Exit(0)
	}

	// Ensure the wallet exists or create it when the create flag is set.
	netDir := networkDir(cfg.AppDataDir.Value, activeNet.Params)
	dbPath := filepath.Join(netDir, walletDbName)

	if cfg.CreateTemp && cfg.Create {
		err := fmt.Errorf("The flags --create and --createtemp can not " +
			"be specified together. Use --help for more information.")
		fmt.Fprintln(os.Stderr, err)
		return nil, nil, err
	}

	dbFileExists, err := cfgutil.FileExists(dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, nil, err
	}

	if cfg.CreateTemp {
		tempWalletExists := false

		if dbFileExists {
			str := fmt.Sprintf("The wallet already exists. Loading this " +
				"wallet instead.")
			fmt.Fprintln(os.Stdout, str)
			tempWalletExists = true
		}

		// Ensure the data directory for the network exists.
		if err := checkCreateDir(netDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, nil, err
		}

		if !tempWalletExists {
			// Perform the initial wallet creation wizard.
			if err := createSimulationWallet(&cfg); err != nil {
				fmt.Fprintln(os.Stderr, "Unable to create wallet:", err)
				return nil, nil, err
			}
		}
	} else if cfg.Create {
		// Error if the create flag is set and the wallet already
		// exists.
		if dbFileExists {
			err := fmt.Errorf("The wallet database file `%v` "+
				"already exists.", dbPath)
			fmt.Fprintln(os.Stderr, err)
			return nil, nil, err
		}

		// Ensure the data directory for the network exists.
		if err := checkCreateDir(netDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, nil, err
		}

		// Perform the initial wallet creation wizard.
		if err := createWallet(&cfg); err != nil {
			fmt.Fprintln(os.Stderr, "Unable to create wallet:", err)
			return nil, nil, err
		}

		// Created successfully, so exit now with success.
		os.Exit(0)
	} else if !dbFileExists && !cfg.NoInitialLoad {
		keystorePath := filepath.Join(netDir, keystore.Filename)
		keystoreExists, err := cfgutil.FileExists(keystorePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, nil, err
		}
		if !keystoreExists {
			// err = fmt.Errorf("The wallet does not exist.  Run with the " +
			// "--create option to initialize and create it...")
			// Ensure the data directory for the network exists.
			fmt.Println("Existing wallet not found in", cfg.ConfigFile.Value)
			if err := checkCreateDir(netDir); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return nil, nil, err
			}

			// Perform the initial wallet creation wizard.
			if err := createWallet(&cfg); err != nil {
				fmt.Fprintln(os.Stderr, "Unable to create wallet:", err)
				return nil, nil, err
			}

			// Created successfully, so exit now with success.
			os.Exit(0)

		} else {
			err = fmt.Errorf("The wallet is in legacy format.  Run with the " +
				"--create option to import it.")
		}
		fmt.Fprintln(os.Stderr, err)
		return nil, nil, err
	}

	// localhostListeners := map[string]struct{}{
	// 	"localhost": {},
	// 	"127.0.0.1": {},
	// 	"::1":       {},
	// }

	// if cfg.UseSPV {
	// 	sac.MaxPeers = cfg.MaxPeers
	// 	sac.BanDuration = cfg.BanDuration
	// 	sac.BanThreshold = cfg.BanThreshold
	// } else {
	if cfg.RPCConnect == "" {
		cfg.RPCConnect = net.JoinHostPort("localhost", activeNet.RPCClientPort)
	}

	// Add default port to connect flag if missing.
	cfg.RPCConnect, err = cfgutil.NormalizeAddress(cfg.RPCConnect,
		activeNet.RPCClientPort)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Invalid rpcconnect network address: %v\n", err)
		return nil, nil, err
	}

	// RPCHost, _, err := net.SplitHostPort(cfg.RPCConnect)
	// if err != nil {
	// 	return nil, nil, err
	// }
	if cfg.EnableClientTLS {
		// if _, ok := localhostListeners[RPCHost]; !ok {
		// 	str := "%s: the --noclienttls option may not be used " +
		// 		"when connecting RPC to non localhost " +
		// 		"addresses: %s"
		// 	err := fmt.Errorf(str, funcName, cfg.RPCConnect)
		// 	fmt.Fprintln(os.Stderr, err)
		// 	fmt.Fprintln(os.Stderr, usageMessage)
		// 	return nil, nil, err
		// }
		// } else {
		// If CAFile is unset, choose either the copy or local pod cert.
		if !cfg.CAFile.ExplicitlySet() {
			cfg.CAFile.Value = filepath.Join(cfg.AppDataDir.Value, defaultCAFilename)

			// If the CA copy does not exist, check if we're connecting to
			// a local pod and switch to its RPC cert if it exists.
			certExists, err := cfgutil.FileExists(cfg.CAFile.Value)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return nil, nil, err
			}
			if !certExists {
				// if _, ok := localhostListeners[RPCHost]; ok {
				podCertExists, err := cfgutil.FileExists(
					podDefaultCAFile)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return nil, nil, err
				}
				if podCertExists {
					cfg.CAFile.Value = podDefaultCAFile
				}
				// }
			}
		}
	}
	// }

	// Only set default RPC listeners when there are no listeners set for
	// the experimental RPC server.  This is required to prevent the old RPC
	// server from sharing listen addresses, since it is impossible to
	// remove defaults from go-flags slice options without assigning
	// specific behavior to a particular string.
	if len(cfg.ExperimentalRPCListeners) == 0 && len(cfg.LegacyRPCListeners) == 0 {
		addrs, err := net.LookupHost("localhost")
		if err != nil {
			return nil, nil, err
		}
		cfg.LegacyRPCListeners = make([]string, 0, len(addrs))
		for _, addr := range addrs {
			addr = net.JoinHostPort(addr, activeNet.RPCServerPort)
			cfg.LegacyRPCListeners = append(cfg.LegacyRPCListeners, addr)
		}
	}

	// Add default port to all rpc listener addresses if needed and remove
	// duplicate addresses.
	cfg.LegacyRPCListeners, err = cfgutil.NormalizeAddresses(
		cfg.LegacyRPCListeners, activeNet.RPCServerPort)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Invalid network address in legacy RPC listeners: %v\n", err)
		return nil, nil, err
	}
	cfg.ExperimentalRPCListeners, err = cfgutil.NormalizeAddresses(
		cfg.ExperimentalRPCListeners, activeNet.RPCServerPort)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Invalid network address in RPC listeners: %v\n", err)
		return nil, nil, err
	}

	// Both RPC servers may not listen on the same interface/port.
	if len(cfg.LegacyRPCListeners) > 0 && len(cfg.ExperimentalRPCListeners) > 0 {
		seenAddresses := make(map[string]struct{}, len(cfg.LegacyRPCListeners))
		for _, addr := range cfg.LegacyRPCListeners {
			seenAddresses[addr] = struct{}{}
		}
		for _, addr := range cfg.ExperimentalRPCListeners {
			_, seen := seenAddresses[addr]
			if seen {
				err := fmt.Errorf("Address `%s` may not be "+
					"used as a listener address for both "+
					"RPC servers", addr)
				fmt.Fprintln(os.Stderr, err)
				return nil, nil, err
			}
		}
	}

	// Only allow server TLS to be disabled if the RPC server is bound to
	// localhost addresses.
	if !cfg.EnableServerTLS {
		allListeners := append(cfg.LegacyRPCListeners,
			cfg.ExperimentalRPCListeners...)
		for _, addr := range allListeners {
			if err != nil {
				str := "%s: RPC listen interface '%s' is " +
					"invalid: %v"
				err := fmt.Errorf(str, funcName, addr, err)
				fmt.Fprintln(os.Stderr, err)
				fmt.Fprintln(os.Stderr, usageMessage)
				return nil, nil, err
			}
			// host, _, err := net.SplitHostPort(addr)
			// if _, ok := localhostListeners[host]; !ok {
			// 	str := "%s: the --noservertls option may not be used " +
			// 		"when binding RPC to non localhost " +
			// 		"addresses: %s"
			// 	err := fmt.Errorf(str, funcName, addr)
			// 	fmt.Fprintln(os.Stderr, err)
			// 	fmt.Fprintln(os.Stderr, usageMessage)
			// 	return nil, nil, err
			// }
		}
	}

	// Expand environment variable and leading ~ for filepaths.
	cfg.CAFile.Value = cleanAndExpandPath(cfg.CAFile.Value)
	cfg.RPCCert.Value = cleanAndExpandPath(cfg.RPCCert.Value)
	cfg.RPCKey.Value = cleanAndExpandPath(cfg.RPCKey.Value)

	// If the pod username or password are unset, use the same auth as for
	// the client.  The two settings were previously shared for pod and
	// client auth, so this avoids breaking backwards compatibility while
	// allowing users to use different auth settings for pod and wallet.
	if cfg.PodUsername == "" {
		cfg.PodUsername = cfg.Username
	}
	if cfg.PodPassword == "" {
		cfg.PodPassword = cfg.Password
	}

	// Warn about missing config file after the final command line parse
	// succeeds.  This prevents the warning on help messages and invalid
	// options.
	if configFileError != nil {
		log.Warnf("%v", configFileError)
	}

	return &cfg, remainingArgs, nil
}

// createDefaultConfig creates a basic config file at the given destination path.
// For this it tries to read the config file for the RPC server (either pod or
// sac), and extract the RPC user and password from it.
func createDefaultConfigFile(destinationPath, serverConfigPath, serverDataDir, walletDataDir string) error {
	// fmt.Println("server config path", serverConfigPath)
	// Read the RPC server config
	serverConfigFile, err := os.Open(serverConfigPath)
	if err != nil {
		return err
	}
	defer serverConfigFile.Close()
	content, err := ioutil.ReadAll(serverConfigFile)
	if err != nil {
		return err
	}
	// content := []byte(samplePodCtlConf)

	// Extract the rpcuser
	rpcUserRegexp, err := regexp.Compile(`(?m)^\s*rpcuser=([^\s]+)`)
	if err != nil {
		return err
	}
	userSubmatches := rpcUserRegexp.FindSubmatch(content)
	if userSubmatches == nil {
		// No user found, nothing to do
		return nil
	}

	// Extract the rpcpass
	rpcPassRegexp, err := regexp.Compile(`(?m)^\s*rpcpass=([^\s]+)`)
	if err != nil {
		return err
	}
	passSubmatches := rpcPassRegexp.FindSubmatch(content)
	if passSubmatches == nil {
		// No password found, nothing to do
		return nil
	}

	// Extract the TLS
	TLSRegexp, err := regexp.Compile(`(?m)^\s*tls=(0|1)(?:\s|$)`)
	if err != nil {
		return err
	}
	TLSSubmatches := TLSRegexp.FindSubmatch(content)

	// Create the destination directory if it does not exists
	err = os.MkdirAll(filepath.Dir(destinationPath), 0700)
	if err != nil {
		return err
	}
	// fmt.Println("config path", destinationPath)
	// Create the destination file and write the rpcuser and rpcpass to it
	dest, err := os.OpenFile(destinationPath,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("ERROR", err)
		return err
	}
	defer dest.Close()

	destString := fmt.Sprintf("username=%s\npassword=%s\n",
		string(userSubmatches[1]), string(passSubmatches[1]))
	if TLSSubmatches != nil {
		fmt.Println("TLS is enabled but more than likely the certificates will fail verification because of the CA. Currently there is no adequate tool for this, but will be soon.")
		destString += fmt.Sprintf("clienttls=%s\n", TLSSubmatches[1])
	}
	output := ";;; Defaults created from local pod/sac configuration:\n" + destString + "\n" + string(sampleModConf)
	dest.WriteString(output)

	return nil
}

func copy(src, dst string) (int64, error) {
	// fmt.Println(src, dst)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
