package db

import (
	"log"
	"os"
	"testing"
)

var dbs Databases

func TestDbs(t *testing.T) {
	os.Chdir("../../..")
	dbs = make(Databases)
	dbs.Init()
}

func TestMongoCRUD(t *testing.T) {
	var database Database
	for _, db := range dbs {
		if dbType := (*db).GetConfig().Type; dbType == "mongo" {
			database = &mongoDatabase{cfg: (*db).GetConfig()}
			break
		}
		// fmt.Printf("val: %s\n", (*db).GetConfig().Name)
	}

	if database != nil {
		database.Init()

		// Select
		rows := database.Select("Sakila_films", "{\"_id\":{\"$lt\": 3}}")
		log.Println(len(rows))

		// Insert
		row := map[string]interface{}{
			"Title":       "test1",
			"Description": "testdesc",
			"ReleaseYear": 2010,
			"Length":      90,
			"ext_id":      1001,
		}
		inDto := InsertDto{
			"Sakila_films",
			"_id",
			1,
			row,
		}
		insertErr := database.Insert(inDto)
		if insertErr != nil {
			log.Fatalln(insertErr)
		}

		// Update
		upDto := UpdateDto{
			"Sakila_films",
			"_id",
			6,
			"Rating",
			"test",
		}
		updateErr := database.Update(upDto)
		if updateErr != nil {
			log.Fatalln(updateErr)
		}
	}
}

func TestPostgresCRUD(t *testing.T) {
	var database Database
	for _, db := range dbs {
		// fmt.Printf("val: %s\n", (*db).GetConfig().Name)
		if dbType := (*db).GetConfig().Type; dbType == "postgres" {
			database = &postgresDatabase{cfg: (*db).GetConfig()}
			break
		}
	}

	if database != nil {
		database.Init()

		// Select
		rows := database.Select("film", "film_id > 10 AND film_id < 22")
		log.Println(len(rows))

		// Insert
		row := map[string]interface{}{
			"title":        "test1",
			"description":  "testdesc",
			"release_year": 2010,
			"length":       90,
			"language_id":  2,
		}
		inDto := InsertDto{
			"Sakila_films",
			"_id",
			1,
			row,
		}
		insertErr := database.Insert(inDto)
		if insertErr != nil {
			log.Fatalln(insertErr)
		}

		// Update
		upDto := UpdateDto{
			"Sakila_films",
			"_id",
			6,
			"Rating",
			"test",
		}
		updateErr := database.Update(upDto)
		if updateErr != nil {
			log.Fatalln(updateErr)
		}
	}
}
