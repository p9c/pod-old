package ctl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.parallelcoin.io/dev/pod/pkg/pod"
	"git.parallelcoin.io/dev/pod/pkg/rpc/json"
	"git.parallelcoin.io/dev/pod/pkg/util"
)

// unusableFlags are the command usage flags which this utility are not able to use.  In particular it doesn't support websockets and consequently notifications.
const unusableFlags = json.UFWebsocketOnly | json.UFNotification

var DefaultConfigFile = filepath.Join(PodCtlHomeDir, "conf.json")

var DefaultRPCCertFile = filepath.Join(NodeHomeDir, "rpc.cert")

var DefaultRPCServer = "127.0.0.1:11048"

var DefaultWallet = "127.0.0.1:11046"

var DefaultWalletCertFile = filepath.Join(SPVHomeDir, "rpc.cert")

var NodeHomeDir = util.AppDataDir("pod", false)

var PodCtlHomeDir = util.AppDataDir("pod/ctl", false)

var SPVHomeDir = util.AppDataDir("pod/spv", false)

// ListCommands categorizes and lists all of the usable commands along with their one-line usage.
func ListCommands() {

	const (
		categoryChain uint8 = iota
		categoryWallet
		numCategories
	)

	// Get a list of registered commands and categorize and filter them.
	cmdMethods := json.RegisteredCmdMethods()
	categorized := make([][]string, numCategories)

	for _, method := range cmdMethods {

		flags, err := json.MethodUsageFlags(method)

		if err != nil {

			// This should never happen since the method was just returned from the package, but be safe.
			continue
		}

		// Skip the commands that aren't usable from this utility.

		if flags&unusableFlags != 0 {

			continue
		}

		usage, err := json.MethodUsageText(method)

		if err != nil {

			// This should never happen since the method was just returned from the package, but be safe.
			continue
		}

		// Categorize the command based on the usage flags.
		category := categoryChain

		if flags&json.UFWalletOnly != 0 {

			category = categoryWallet
		}

		categorized[category] = append(categorized[category], usage)
	}

	// Display the command according to their categories.
	categoryTitles := make([]string, numCategories)
	categoryTitles[categoryChain] = "Chain Server Commands:"
	categoryTitles[categoryWallet] = "Wallet Server Commands (--wallet):"

	for category := uint8(0); category < numCategories; category++ {

		fmt.Println(categoryTitles[category])
		fmt.Println()

		for _, usage := range categorized[category] {

			fmt.Println("  ", usage)
		}

		fmt.Println()
	}

}

// cleanAndExpandPath expands environement variables and leading ~ in the passed path, cleans the result, and returns it.
func cleanAndExpandPath(
	path string,

) string {

	// Expand initial ~ to OS specific home directory.
	if strings.HasPrefix(path, "~") {

		homeDir := filepath.Dir(PodCtlHomeDir)
		path = strings.Replace(path, "~", homeDir, 1)
	}

	// NOTE: The os.ExpandEnv doesn't work with Windows-style %VARIABLE%, but they variables can still be expanded via POSIX-style $VARIABLE.
	return filepath.Clean(os.ExpandEnv(path))
}

// loadConfig initializes and parses the config using a config file and command line options.
// The configuration proceeds as follows:
// 	1) Start with a default config with sane settings
// 	2) Pre-parse the command line to check for an alternative config file
// 	3) Load configuration file overwriting defaults with any specified options
// 	4) Parse CLI options and overwrite/add any specified options
// The above results in functioning properly without any config settings while still allowing the user to override settings with config files and command line options.  Command line options always take precedence.
func loadConfig() error {
	// Default config.
	cfg := pod.Config{
		ConfigFile: &DefaultConfigFile,
		RPCConnect: &DefaultRPCServer,
		RPCCert:    &DefaultRPCCertFile,
	}
	// Override the RPC certificate if the --wallet flag was specified and the user did not specify one.
	if *cfg.Wallet && *cfg.RPCCert == DefaultRPCCertFile {
		*cfg.RPCCert = DefaultWalletCertFile
	}
	return nil
}
