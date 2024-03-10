package main

import (
	"log"
	"os"

	"github.com/fernandoocampo/basic-micro/internal/application"
)

func main() {
	log.Println("starting application")

	if err := application.New().Run(); err != nil {
		log.Printf("unable to start service: %s", err)
		os.Exit(-1)
	}

	log.Println("finishing application")
}
