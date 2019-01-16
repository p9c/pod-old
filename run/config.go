package pod

import (
	"fmt"
	"os"
	"runtime"

	flags "github.com/jessevdk/go-flags"
)

var (
	cfg *Config
)

// LoadConfig loads configuration by loading default, overlaying config file settings and then finally overriding with command line parameters
func LoadConfig() (config *Config, args []string, err error) {
	cfg = &Config{}
	config = cfg
	args = os.Args
	serviceOpts := serviceOptions{}
	parser := flags.NewParser(cfg, flags.HelpFlag)
	if runtime.GOOS == "windows" {
		parser.AddGroup("Service Options", "Service Options", &serviceOpts)
	}
	_, err = parser.Parse()
	if err != nil {
		// fmt.Println("pod", version())
		fmt.Println(err.Error())
		return
	}
	return
}

type serviceOptions struct {
	ServiceCommand string `short:"s" long:"service" description:"Service command {install, remove, start, stop}"`
}
