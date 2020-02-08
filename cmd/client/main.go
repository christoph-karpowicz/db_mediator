package main

import (
	"log"
	"os"

	"github.com/christoph-karpowicz/unifier/internal/client/application"

	_ "github.com/lib/pq"
)

func main() {

	var app application.Application = application.Application{}

	app.Init()

	app.SetCLI()

	err := app.CLI.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
