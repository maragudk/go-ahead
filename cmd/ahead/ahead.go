package main

import (
	"flag"
	"log"

	"go-ahead"
)

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
		ExternalPort: c.ExternalPort,
		InternalPort: c.InternalPort,
	})

	if err := s.Start(); err != nil {
		log.Fatalln("Could not start:", err)
	}
}
