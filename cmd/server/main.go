package main

import (
	"flag"
	"log"

	"go-ahead/server"
	"go-ahead/storage"
)

func main() {
	// Disable timestamps in logger
	log.SetFlags(0)

	path := flag.String("config", "development.toml", "Path of config file in TOML format")
	flag.Parse()
	c, err := ReadConfig(*path)
	if err != nil {
		log.Fatalln("Could not read config file:", err)
	}

	s := server.New(server.Options{
		Storer:       createStorer(c),
		ExternalPort: c.ExternalPort,
		InternalPort: c.InternalPort,
		Name:         c.Name,
		Version:      Version,
	})

	if err := s.Start(); err != nil {
		log.Fatalln("Could not start:", err)
	}
}

func createStorer(c Config) *storage.Storer {
	return storage.NewStorer(storage.NewStorerOptions{
		User:     c.Storer.User,
		Host:     c.Storer.Host,
		Port:     c.Storer.Port,
		Database: c.Storer.Database,
		Cert:     c.Storer.Cert,
		Key:      c.Storer.Key,
		RootCert: c.Storer.RootCert,
	})
}
