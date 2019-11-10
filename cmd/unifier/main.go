package main

import (
	"fmt"
	"log"
	"os"

	"github.com/christoph-karpowicz/unifier/internal/application"

	_ "github.com/lib/pq"
)

func main() {

	var app application.Application = application.Application{}

	app.SetCLI()

	err := app.CLI.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	app.Init()

	fmt.Println(app.Lang)

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
