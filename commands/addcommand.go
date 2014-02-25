package commands

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type AddCommand struct {
	Source      *url.URL
	Destination string
	IsForce     bool
	IsDownload  bool
}

func initAddFlags(addCommand *AddCommand) *flag.FlagSet {
	addFlags := flag.NewFlagSet("add", flag.ExitOnError)

	addFlags.BoolVar(&addCommand.IsForce, "force", false, "Use -force to override existing values")
	addFlags.BoolVar(&addCommand.IsForce, "f", false, "Use -f to override existing values")
	addFlags.BoolVar(&addCommand.IsDownload, "download", false, "Use -download to download the source to the destination")
	addFlags.BoolVar(&addCommand.IsDownload, "dl", false, "Use -dl to download the source to the destination")

	return addFlags
}

func ParseAddCommand(flags []string) (Command, error) {
	addCommand := AddCommand{}

	addFlags := initAddFlags(&addCommand)

	addFlags.Parse(flags)

	if addFlags.NArg() < 2 {
		return nil, errors.New("Not enough params")
	}

	var err error
	addCommand.Source, err = url.Parse(addFlags.Arg(0))

	if err != nil || addCommand.Source.Scheme == "" || addCommand.Source.Host == "" {
		return nil, errors.New(fmt.Sprintf("Can't parse '%s' as url", addFlags.Arg(0)))
	}

	addCommand.Destination = addFlags.Arg(1)

	return &addCommand, nil
}

func downloadAndWriteToFile(url *url.URL, filepath string) error {
	//Open rw, create of non-existing, should not exist already
	out, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	defer out.Close()

	if err != nil {
		if os.IsExist(err) {
			return errors.New(fmt.Sprintf("File '%s' already exist, won't override", filepath))
		} else {
			return errors.New("Error opening file")
		}
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return errors.New("Error downloading file")
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}

func (addCommand *AddCommand) Execute(configLoader ConfigLoader) error {
	config, err := configLoader()

	if err != nil {
		return err
	}

	currentValue, exists := config.Rewrites[addCommand.Source.String()]

	if exists && !addCommand.IsForce {
		log.Printf("Entry { '%s' => '%s' } exists, use option -force to override", addCommand.Source, currentValue)
		return errors.New("Entry exists")
	}

	config.Rewrites[addCommand.Source.String()] = addCommand.Destination

	if addCommand.IsDownload {
		err := downloadAndWriteToFile(addCommand.Source, addCommand.Destination)

		if err != nil {
			log.Println(err)
		}

		return nil
	}

	return nil
}
