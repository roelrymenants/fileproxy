package commands

import "github.com/roelrymenants/fileproxy/proxyconfig"

type Command interface {
	Execute(*proxyconfig.Config) error
}

type CommandChain []Command

func (chain CommandChain) Execute(config *proxyconfig.Config) error {
	for index := range chain {
		cmd := chain[index]
		err := cmd.Execute(config)

		if err != nil {
			return err
		}
	}

	return nil
}
