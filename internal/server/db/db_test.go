package db

import (
	"log"
	"os"
	"testing"
)

var configs *configArray

func TestYAML(t *testing.T) {
	os.Chdir("../../..")
	dbs := Databases{DBMap: make(map[string]*Database)}
	configs = dbs.ImportYAML()
	dbs.ValidateYAML()
}

func TestMongoCRUD(t *testing.T) {
	var database Database
	for i := 0; i < len(configs.Databases); i++ {
		// fmt.Printf("val: %s\n", configs.Databases[i].Name)
		if dbType := configs.Databases[i].Type; dbType == "mongo" {
			database = &mongoDatabase{cfg: &configs.Databases[i]}
			break
		}
	}

	if database != nil {
		database.Init()

		// Select
		rows := database.Select("Sakila_films", "{\"_id\":{\"$lt\": 3}}")
		log.Println(len(rows))

		// Update
		_, err := database.Update("Sakila_films", "_id", 6, "Rating", "test")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func TestPostgresCRUD(t *testing.T) {
	var database Database
	for i := 0; i < len(configs.Databases); i++ {
		// fmt.Printf("val: %s\n", configs.Databases[i].Name)
		if dbType := configs.Databases[i].Type; dbType == "postgres" {
			database = &postgresDatabase{cfg: &configs.Databases[i]}
			break
		}
	}

	if database != nil {
		database.Init()

		// Select
		rows := database.Select("film", "film_id > 10 AND film_id < 22")
		log.Println(len(rows))

		// Update
		_, err := database.Update("film", "film_id", 1, "description", "test")
		if err != nil {
			log.Fatalln(err)
		}
	}
}
