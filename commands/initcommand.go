package commands

import (
	"flag"

	"github.com/roelrymenants/fileproxy/proxyconfig"
)

type InitCommand struct {
	IsVerbose bool
}

func ParseInitCommand(args []string) (Command, error) {
	initCmd := InitCommand{}
	
	initFlags := flag.NewFlagSet("init", flag.ExitOnError)

	initFlags.BoolVar(&initCmd.IsVerbose, "verbose", false, "Use -verbose to make the proxy verbose")
	initFlags.BoolVar(&initCmd.IsVerbose, "v", false, "Use -v to make the proxy verbose")
	
	initFlags.Parse(args)

	return &InitCommand{}, nil
}

//Config should alwyas be nil
func (initCmd *InitCommand) Execute(configLoader ConfigLoader) error {
	config := proxyconfig.NewConfig()
	config.Verbose = initCmd.IsVerbose
	config.SaveToFile(proxyconfig.DefaultConfigFile)

	configLoader()

	return nil
}
