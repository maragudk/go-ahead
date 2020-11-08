package main

import (
	"github.com/BurntSushi/toml"

	"github.com/maragudk/go-ahead/errors2"
)

// Config holds configuration data read from a config file.
type Config struct {
	Host    string
	Port    int
	Cert    string
	Key     string
	Emailer struct {
		Token string
	}
	Database struct {
		User     string
		Password string
		Host     string
		Port     int
		Socket   string
		Name     string
	}
}

// ReadConfig from path.
func ReadConfig(path string) (Config, error) {
	var c Config
	_, err := toml.DecodeFile(path, &c)
	return c, errors2.Wrap(err, "could not read config file at %v", path)
}
