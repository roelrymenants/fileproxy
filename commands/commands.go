package commands

import "github.com/roelrymenants/fileproxy/proxyconfig"

type ConfigLoader func() (*proxyconfig.Config, error)

type Command interface {
	Execute(ConfigLoader) error
}

type CommandChain []Command

func (chain CommandChain) Execute(configLoader ConfigLoader) error {
	for index := range chain {
		cmd := chain[index]
		err := cmd.Execute(configLoader)

		if err != nil {
			return err
		}
	}

	return nil
}
