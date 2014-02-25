package commands

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/roelrymenants/fileproxy/proxyconfig"
)

type RemoveCommand struct {
	Source *url.URL
}

func ParseRemoveCommand(args []string) (*RemoveCommand, error) {
	removeCommand := &RemoveCommand{}

	var err error
	removeCommand.Source, err = url.Parse(args[0])

	if err != nil || removeCommand.Source.Scheme == "" || removeCommand.Source.Host == "" {
		return nil, errors.New(fmt.Sprintf("Can't parse '%s' as url", args[0]))
	}

	return removeCommand, nil
}

func (removeCommand *RemoveCommand) Execute(config *proxyconfig.Config) error {
	key := removeCommand.Source.String()
	_, ok := config.Rewrites[key]

	if !ok {
		return errors.New(fmt.Sprintf("No such key '%s' in redirects", key))
	}

	delete(config.Rewrites, key)

	return nil
}
