package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/christoph-karpowicz/unifier/internal/db"

	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

var app = cli.NewApp()

func info() {
	app.Name = "Unifier CLI"
	app.Usage = "Database synchronization app."
	app.Author = "Krzysztof Karpowicz"
	app.Version = "1.0.0"
}

func commands() {
	app.Commands = []cli.Command{
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
}

func flags(lang *string) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "lang",
			Value:       "english",
			Usage:       "language for the greeting",
			Destination: lang,
		},
	}
}

func main() {

	var lang string

	info()
	commands()
	flags(&lang)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	databasesConfigFile, err := os.Open("config/databases.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened databases.json")
	defer databasesConfigFile.Close()

	// Read opened file as a byte array.
	byteValue, _ := ioutil.ReadAll(databasesConfigFile)

	var databases db.Databases

	json.Unmarshal(byteValue, &databases)

	for i := 0; i < len(databases.Databases); i++ {
		fmt.Println("Database Type: " + databases.Databases[i].Type)
		fmt.Println("Database Name: " + databases.Databases[i].Name)
	}

	fmt.Println(lang)

	// fmt.Printf("%v\n\n", os.Args)

	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbname)
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Successfully connected!")

	// rows, err := db.Query(`SELECT film_id, title FROM film WHERE title ILIKE 'Des%'`)
	// for rows.Next() {
	// 	var id string
	// 	var title string

	// 	if err := rows.Scan(&id, &title); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("id: %s, title: %s\n", id, title)
	// }

}
