package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/roelrymenants/fileproxy/commands"
	"github.com/roelrymenants/fileproxy/proxyconfig"
)

type CommandParser func(flags []string) (commands.Command, error)

var CommandParsers = map[string]CommandParser{
	"init":   commands.ParseInitCommand,
	"add":    commands.ParseAddCommand,
	"remove": commands.ParseRemoveCommand,
	"run":    commands.ParseRunCommand,
}

func fetchCommand(args []string) (commands.Command, error) {
	cmdFunc, ok := CommandParsers[args[1]]

	if !ok {
		return nil, errors.New("Command not found")
	}

	return cmdFunc(args[2:])
}

func main() {
	fSet := flag.NewFlagSet("test", flag.ExitOnError)
	fSet.Bool("test", false, "testflag")

	flag.Parse()

	flag.Usage = func() {
		fmt.Printf("Usage: %s init|add|remove|run [options]\n", os.Args[0])
		fmt.Printf("Try %s <cmd> -help for options\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	cmd, err := fetchCommand(os.Args)

	if err != nil {
		log.Printf("Error parsing command: %s", err)
		flag.Usage()
		return
	}

	var config *proxyconfig.Config

	configLoader := func() (*proxyconfig.Config, error) {
		if config == nil {
			var err error //Can't use lazy init cause we'd override the closure
			config, err = proxyconfig.LoadConfig(proxyconfig.DefaultConfigFile)

			if err != nil {
				log.Fatalf("Could not load config file '%s'. Try running '%s init'.", proxyconfig.DefaultConfigFile, os.Args[0])
			}
		}

		return config, nil
	}

	err = cmd.Execute(configLoader)

	if err == nil {
		//All went well
		config.SaveToFile(proxyconfig.DefaultConfigFile)
		log.Printf("%+v", config)
	} else {
		log.Printf("%s", err)
	}
}
