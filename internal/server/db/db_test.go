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
	// fmt.Println(configs)
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
		_, err := database.Update("Sakila_films", "_id", 6, "Rating", "test2222")
		if err != nil {
			log.Println(err)
		}
	}
}
