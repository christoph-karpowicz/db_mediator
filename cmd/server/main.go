package main

import (
	"github.com/christoph-karpowicz/unifier/internal/server/application"

	_ "github.com/lib/pq"
)

func main() {

	var App application.Application = application.Application{}
	App.Init()

}
