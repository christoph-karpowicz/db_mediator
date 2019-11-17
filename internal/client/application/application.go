package application

import (
	"fmt"

	"github.com/urfave/cli"
)

type Application struct {
	CLI  *cli.App
	Lang string
}

func (a *Application) Init() {

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
				fmt.Println("one-off")
				// a.synchronizeArray(c.Args())
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
