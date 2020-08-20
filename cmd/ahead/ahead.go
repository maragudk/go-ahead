package main

import (
	"flag"
	"log"

	"go-ahead"
	"go-ahead/storage"
)

const Version = "VERSION"

func main() {
	// Disable timestamps in logger
	log.SetFlags(0)

	path := flag.String("config", "ahead.toml", "Path of config file in TOML format")
	flag.Parse()
	c, err := ReadConfig(*path)
	if err != nil {
		log.Fatalln("Could not read config file:", err)
	}

	s := ahead.NewServer(ahead.NewServerOptions{
		Storer:       createStorer(c),
		ExternalPort: c.ExternalPort,
		InternalPort: c.InternalPort,
		Version:      Version,
	})

	if err := s.Start(); err != nil {
		log.Fatalln("Could not start:", err)
	}
}

func createStorer(c Config) *storage.Storer {
	return storage.NewStorer(storage.NewStorerOptions{
		User:     c.Database.User,
		Host:     c.Database.Host,
		Port:     c.Database.Port,
		Database: c.Database.Database,
		Cert:     c.Database.Cert,
		Key:      c.Database.Key,
		RootCert: c.Database.RootCert,
	})
}
