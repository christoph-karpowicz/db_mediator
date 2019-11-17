package application

import (
	"fmt"

	"github.com/christoph-karpowicz/unifier/internal/db"
	"github.com/christoph-karpowicz/unifier/internal/synch"

	"github.com/urfave/cli"
)

type Application struct {
	CLI    *cli.App
	Lang   string
	dbs    *db.Databases
	synchs *synch.Synchs
}

func (a *Application) Init() {
	a.dbs = &db.Databases{DBMap: make(map[string]*db.Database)}
	a.dbs.ImportJSON()
	a.dbs.ValidateJSON()
	a.synchs = &synch.Synchs{SynchMap: make(map[string]*synch.Synch)}
	a.synchs.ImportJSON()
	a.synchs.ValidateJSON()
}

func (a *Application) SetCLI() {
	a.CLI = cli.NewApp()
	a.CLI.Name = "Unifier CLI"
	a.CLI.Usage = "Database synchronization app."
	a.CLI.Author = "Krzysztof Karpowicz"
	a.CLI.Version = "1.0.0"

	a.CLI.Commands = []cli.Command{
		{
			Name:    "one-off",
			Aliases: []string{"oo"},
			Usage:   "One off synchronization.",
			Action: func(c *cli.Context) {
				a.synchronizeArray(c.Args())
			},
		},
		{
			Name:    "ongoing",
			Aliases: []string{"ng"},
			Usage:   "Start ongoing synchronization.",
			Action: func(c *cli.Context) {
				fmt.Println("ong")
			},
		},
	}

	a.CLI.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "lang",
			Value:       "english",
			Usage:       "language for the greeting",
			Destination: &a.Lang,
		},
	}
}

func (a *Application) synchronize(synchKey string) {
	synch := a.synchs.SynchMap[synchKey]
	if synch == nil {
		panic("Synch '" + synchKey + "' not found.")
	}

	synch.SetDatabases(a.dbs.DBMap)

	fmt.Println(*synch)
}

func (a *Application) synchronizeArray(synchKeys []string) {
	for _, arg := range synchKeys {
		a.synchronize(arg)
	}
}
