package fileproxy

import (
	"io/ioutil"
	"log"

	"encoding/json"

	"github.com/howeyc/fsnotify"
)

const DefaultConfigFile string = "rewrites.json"

type Config struct {
	Rewrites map[string]string
}

func NewConfig() *Config {
	return &Config{Rewrites: make(map[string]string)}
}

func LoadConfig(filepath string) *Config {
	var config Config
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(b, &config)
	log.Printf("%v", config)

	return &config
}

func (config *Config) SaveToFile(filepath string) {
	serialConfig, err := json.Marshal(config)

	if err != nil {
		log.Fatalf("Error (%s) writing config file", err)
	}

	ioutil.WriteFile(filepath, serialConfig, 0644)
}

func StartWatching(filepath string) chan Config {
	configChan := make(chan Config)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				log.Printf("Event: %+v", event)
				config := LoadConfig(filepath)

				if config != nil {
					configChan <- *config
				}
			case err := <-watcher.Error:
				log.Panic("Error: %+v", err)
			}
		}
	}()

	err = watcher.Watch(filepath)
	if err != nil {
		log.Fatal("Error opening watcher", err)
	}

	return configChan
}
