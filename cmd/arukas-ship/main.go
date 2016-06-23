package main

import (
	"github.com/yamamoto-febc/arukas-ship"
	"log"
)

func main() {

	config, err := ship.InitializeConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting arukas-ship server on: 0.0.0.0:%d", config.Serve.Port)
	if err := ship.Serve(config); err != nil {
		log.Fatal(err)
	}
}
