package main

import (
	"github.com/BurntSushi/toml"

	"go-ahead/errors2"
)

type Config struct {
	Name         string
	ExternalPort int
	InternalPort int
	Emailer      struct {
		Token string
	}
	Storer struct {
		Host     string
		Port     int
		User     string
		Database string
		Cert     string
		Key      string
		RootCert string
	}
}

// ReadConfig from path.
func ReadConfig(path string) (Config, error) {
	var c Config
	_, err := toml.DecodeFile(path, &c)
	return c, errors2.Wrap(err, "could not read config file at %v", path)
}
