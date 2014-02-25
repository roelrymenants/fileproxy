package commands

import "github.com/roelrymenants/fileproxy"

type Command interface {
	Execute(*fileproxy.Config) error
}

type CommandChain []Command

func (chain CommandChain) Execute(config *fileproxy.Config) error {
	for index := range chain {
		cmd := chain[index]
		err := cmd.Execute(config)

		if err != nil {
			return err
		}
	}

	return nil
}
