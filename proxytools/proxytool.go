package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/roelrymenants/fileproxy"
	"github.com/roelrymenants/proxytool/commands"
)

func main() {
	var cmd commands.Command

	fSet := flag.NewFlagSet("test", flag.ExitOnError)
	fSet.Bool("test", false, "testflag")

	flag.Parse()

	flag.Usage = func() {
		fmt.Printf("Usage: %s add|remove [options]\n", os.Args[0])
		fmt.Printf("Try %s <cmd> -help for options\n", os.Args[0])
	}

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}
	
	var err error

	switch os.Args[1] {
	case "run":
		cmd, err = commands.(os.Args[2:])
	case "add":
		cmd, err = commands.ParseAddCommand(os.Args[2:])
	case "remove":
		cmd, err = commands.ParseRemoveCommand(os.Args[2:])
	/*case "replace":
	addCmd := commands.ParseAddCommand(os.Args[2:])
	removeCmd := &commands.RemoveCommand{addCmd.Source}

	cmd = commands.CommandChain([]commands.Command{removeCmd, addCmd})*/
	case "init":
		//Special case, no actual command
		config := fileproxy.NewConfig()
		config.SaveToFile(fileproxy.DefaultConfigFile)

		return
	}
	
	if err != nil {
		log.Printf("Error parsing command: %s", err)
		flag.Usage()
		return
	}

	config, err := fileproxy.LoadConfig(fileproxy.DefaultConfigFile)

	if err != nil {
		log.Fatalf("Could not load config file '%s", fileproxy.DefaultConfigFile)
	}

	err = cmd.Execute(config)

	if err == nil {
		//All went well
		config.SaveToFile(fileproxy.DefaultConfigFile)
		log.Printf("%+v", config)
	} else {
		log.Printf("%s", err)
	}
}
