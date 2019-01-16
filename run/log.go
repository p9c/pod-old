package pod

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jrick/logrotate/rotator"
)

// logWriter implements an io.Writer that outputs to both standard output and the write-end pipe of an initialized log rotator.
type logWriter struct{}

func (logWriter) Write(p []byte) (n int, err error) {
	os.Stdout.Write(p)
	logRotator.Write(p)
	return len(p), nil
}

// Loggers per subsystem.  A single backend logger is created and all subsytem loggers created from it will write to the backend.  When adding new subsystems, add the subsystem logger variable here and to the  subsystemLoggers map. Loggers can not be used before the log rotator has been initialized with a log file.  This must be performed early during application startup by calling initLogRotator.
var (
	// logRotator is one of the logging outputs.  It should be closed on application shutdown.
	logRotator *rotator.Rotator
)

// Initialize package-global logger variables.
func init() {
}

// initLogRotator initializes the logging rotater to write logs to logFile and create roll files in the same directory.  It must be called before the package-global log rotater variables are used.
func initLogRotator(logFile string) {
	logDir, _ := filepath.Split(logFile)
	err := os.MkdirAll(logDir, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create log directory: %v\n", err)
		os.Exit(1)
	}
	r, err := rotator.New(logFile, 10*1024, false, 3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create file rotator: %v\n", err)
		os.Exit(1)
	}
	logRotator = r
}

// // setLogLevel sets the logging level for provided subsystem.  Invalid subsystems are ignored.  Uninitialized subsystems are dynamically created as needed.
// func setLogLevel(subsystemID string, logLevel string) {
// 	// Ignore invalid subsystems.
// 	logger, ok := subsystemLoggers[subsystemID]
// 	if !ok {
// 		return
// 	}
// 	// Defaults to info if the log level is invalid.
// 	level, _ := log.LevelFromString(logLevel)
// 	logger.SetLevel(level)
// }

// // setLogLevels sets the log level for all subsystem loggers to the passed level.  It also dynamically creates the subsystem loggers as needed, so it can be used to initialize the logging system.
// func setLogLevels(logLevel string) {
// 	// Configure all sub-systems with the new logging level.  Dynamically create loggers as needed.
// 	for subsystemID := range subsystemLoggers {
// 		setLogLevel(subsystemID, logLevel)
// 	}
// }

// directionString is a helper function that returns a string that represents the direction of a connection (inbound or outbound).
func directionString(inbound bool) string {
	if inbound {
		return "inbound"
	}
	return "outbound"
}

// pickNoun returns the singular or plural form of a noun depending on the count n.
func pickNoun(n uint64, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}
