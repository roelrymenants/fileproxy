package main

import (
	"testing"

	"github.com/roelrymenants/fileproxy"
	"github.com/roelrymenants/proxytool/commands"
)

var source, destination, firstDestination = "http://www.google.com/some/source", "/some/destination", "/first/destination"
var completeAddExample = []string{"proxytool", "add", "-force", source, destination}

func TestAddCommandFlagParsing(t *testing.T) {
	addCommand := commands.ParseAddCommand(completeAddExample[2:])

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
	addCommand := commands.ParseAddCommand(completeAddExample[2:])
	addCommand.IsForce = false

	config := fileproxy.NewConfig()

	addCommand.Execute(config)

	put, ok := config.Rewrites[source]

	if !ok || put != destination {
		t.Errorf("Source-dest mapping not made")
	}
}

func TestAddCommandForceOverwrite(t *testing.T) {
	addCommand := commands.ParseAddCommand(completeAddExample[2:])

	config := fileproxy.NewConfig()
	config.Rewrites[source] = firstDestination

	addCommand.Execute(config)

	put, ok := config.Rewrites[source]

	if !ok || put != destination {
		t.Errorf("Source-dest mapping not made")
	}
}

func TestAddCommandNoForceNoOverwrite(t *testing.T) {
	addCommand := commands.ParseAddCommand(completeAddExample[2:])
	addCommand.IsForce = false

	config := fileproxy.NewConfig()
	config.Rewrites[source] = firstDestination

	err := addCommand.Execute(config)

	put, ok := config.Rewrites[source]

	if err == nil {
		t.Errorf("No error thrown on existing key")
	}

	if !ok || put == destination {
		t.Errorf("Source-dest mapping made, force not defined")
	}
}
