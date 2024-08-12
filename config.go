package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Place initConfig() in config.go
func initConfig() Config {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}
	return config
}
