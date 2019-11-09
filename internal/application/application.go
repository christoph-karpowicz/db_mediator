package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/christoph-karpowicz/unifier/internal/db"

	"github.com/urfave/cli"
)

type Application struct {
	CLI  *cli.App
	Lang string
	dbs  db.Databases
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

func (a *Application) ImportDatabases() {
	databasesConfigFile, err := os.Open("config/databases.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Parsing databases.json...")
	defer databasesConfigFile.Close()

	// Read opened file as a byte array.
	byteValue, _ := ioutil.ReadAll(databasesConfigFile)

	json.Unmarshal(byteValue, &a.dbs)

	fmt.Println("----------------")
	fmt.Println("Databases:")
	// for i := 0; i < len(a.dbs.Databases); i++ {
	// 	fmt.Println("	- Name: " + a.dbs.Databases[i].Name + ", type: " + a.dbs.Databases[i].Type)
	// }
	fmt.Println(a.dbs.Databases)
	fmt.Println(a.dbs.Databases["dvdrental"].Type)
	fmt.Println("----------------")
}
