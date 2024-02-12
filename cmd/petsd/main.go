package main

import (
	"log"
	"os"

	"github.com/fernandoocampo/basic-micro/internal/application"
)

var (
	Version    string
	BuildDate  string
	CommitHash string
)

func main() {
	log.Println("starting application")

	app := newApplicationServer()

	if err := app.Run(); err != nil {
		log.Printf("unable to start service: %s", err)
		os.Exit(-1)
	}

	log.Println("finishing application")
}

func newApplicationServer() *application.Server {
	settings := application.Setup{
		Version:    Version,
		BuildDate:  BuildDate,
		CommitHash: CommitHash,
	}

	return application.NewServer(settings)
}
