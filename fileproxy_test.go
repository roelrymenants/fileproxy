package main

import (
	"testing"

	"github.com/roelrymenants/fileproxy/commands"
	"github.com/roelrymenants/fileproxy/proxyconfig"
)

var source, destination, firstDestination = "http://www.google.com/some/source", "/some/destination", "/first/destination"
var completeAddExample = []string{"proxytool", "add", "-force", source, destination}

func TestAddCommandFlagParsing(t *testing.T) {
	cmd, err := commands.ParseAddCommand(completeAddExample[2:])

	if err != nil {
		t.Error("Error parsing command")
	}

	addCommand, ok := cmd.(*commands.AddCommand)

	if !ok {
		t.Fatal("Parsed command is not an AddCommand")
	}

	if !addCommand.IsForce {
		t.Error("Force option not parsed")
	}

	if addCommand.Source.String() != source {
		t.Errorf("Source not parsed, found %s", addCommand.Source)
	}

	if addCommand.Destination != destination {
		t.Errorf("Destination not parsed, found %s", addCommand.Destination)
	}
}

func TestAddCommandBasic(t *testing.T) {
	cmd, err := commands.ParseAddCommand(completeAddExample[2:])

	if err != nil {
		t.Error("Error parsing command")
	}

	addCommand, ok := cmd.(*commands.AddCommand)

	if !ok {
		t.Fatal("Parsed command is not an AddCommand")
	}

	addCommand.IsForce = false
	config := proxyconfig.NewConfig()

	addCommand.Execute(func() (*proxyconfig.Config, error) {
		return config, nil
	})
	
	put, ok := config.Rewrites[source]	

	if !ok || put != destination {
		t.Errorf("Source-dest mapping not made")
	}
}

func TestAddCommandForceOverwrite(t *testing.T) {
	cmd, err := commands.ParseAddCommand(completeAddExample[2:])

	if err != nil {
		t.Error("Error parsing command")
	}

	addCommand, ok := cmd.(*commands.AddCommand)

	if !ok {
		t.Fatal("Parsed command is not an AddCommand")
	}

	config := proxyconfig.NewConfig()
	config.Rewrites[source] = firstDestination

	addCommand.Execute(func() (*proxyconfig.Config, error) {
		return config, nil
	})

	put, ok := config.Rewrites[source]

	if !ok || put != destination {
		t.Errorf("Source-dest mapping not made")
	}
}

func TestAddCommandNoForceNoOverwrite(t *testing.T) {
	cmd, err := commands.ParseAddCommand(completeAddExample[2:])

	if err != nil {
		t.Error("Error parsing command")
	}

	addCommand, ok := cmd.(*commands.AddCommand)

	if !ok {
		t.Fatal("Parsed command is not an AddCommand")
	}

	addCommand.IsForce = false

	config := proxyconfig.NewConfig()
	config.Rewrites[source] = firstDestination

	err = addCommand.Execute(func() (*proxyconfig.Config, error) {
		return config, nil
	})

	put, ok := config.Rewrites[source]

	if err == nil {
		t.Errorf("No error thrown on existing key")
	}

	if !ok || put == destination {
		t.Errorf("Source-dest mapping made, force not defined")
	}
}
