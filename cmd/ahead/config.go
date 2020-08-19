package main

import (
	"github.com/BurntSushi/toml"

	"go-ahead/errors2"
)

type Config struct {
	ExternalPort int
	InternalPort int
}

// ReadConfig from path.
func ReadConfig(path string) (Config, error) {
	var c Config
	_, err := toml.DecodeFile(path, &c)
	return c, errors2.Wrap(err, "could not read config file at %v", path)
}