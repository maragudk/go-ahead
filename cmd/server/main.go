// Package main is the entry point to the application.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/maragudk/go-ahead/comms"
	"github.com/maragudk/go-ahead/server"
	"github.com/maragudk/go-ahead/storage"
)

func main() {
	logger := log.New(os.Stdout, "", 0)

	path := flag.String("config", "development.toml", "Path of config file in TOML format")
	flag.Parse()
	c, err := ReadConfig(*path)
	if err != nil {
		logger.Fatalln("Could not read config file:", err)
	}

	s := server.New(server.Options{
		Emailer: createEmailer(c),
		Storer:  createStorer(c, logger),
		Logger:  logger,
		Host:    c.Host,
		Port:    c.Port,
		Cert:    c.Cert,
		Key:     c.Key,
	})

	if err := s.Start(); err != nil {
		logger.Fatalln("Could not start:", err)
	}
}

func createStorer(c Config, logger *log.Logger) *storage.Storer {
	return storage.New(storage.Options{
		Logger:   logger,
		User:     c.Database.User,
		Password: c.Database.Password,
		Host:     c.Database.Host,
		Port:     c.Database.Port,
		Socket:   c.Database.Socket,
		Name:     c.Database.Name,
	})
}

func createEmailer(c Config) *comms.Emailer {
	return comms.NewEmailer(c.Emailer.Token)
}
