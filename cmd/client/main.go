package main

import (
	"github.com/christoph-karpowicz/unifier/internal/client/application"

	_ "github.com/lib/pq"
)

func main() {

	var app application.Application = application.Application{}

	app.Init()
}
