package main

import (
	"log"

	"go-ahead"
)

func main() {
	if err := ahead.Start(ahead.StartOptions{}); err != nil {
		log.Fatalln("Could not start:", err)
	}
}
