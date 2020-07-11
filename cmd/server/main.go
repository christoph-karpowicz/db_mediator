package main

import (
	"os"

	"github.com/christoph-karpowicz/unifier/internal/server/application"

	_ "github.com/lib/pq"
)

func main() {

	args := os.Args
	if len(args) > 1 && args[1] == "debug" {
		os.Chdir("../..")
	}

	var App application.Application = application.Application{}
	App.Init()

}
