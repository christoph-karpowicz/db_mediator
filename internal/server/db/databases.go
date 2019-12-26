/*
Package db contains database configurations and
methods for querying.
*/
package db

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

// Databases imports, validates and holds information about databases from JSON config files.
type Databases struct {
	DBMap map[string]*Database
}

// ImportYAML parses and saves YAML config files.
func (d *Databases) ImportYAML() {
	databasesFilePath, _ := filepath.Abs("./config/databases.yaml")

	databasesConfigFile, err := os.Open(databasesFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsing %s...\n", databasesFilePath)
	defer databasesConfigFile.Close()

	byteArray, err := ioutil.ReadAll(databasesConfigFile)
	if err != nil {
		panic(err)
	}

	var dbDataArr databaseDataArray = databaseDataArray{}
	marshalErr := yaml.Unmarshal(byteArray, &dbDataArr)
	if marshalErr != nil {
		log.Fatalf("error: %v", marshalErr)
	}

	fmt.Println("----------------")
	fmt.Println("Databases:")
	for i := 0; i < len(dbDataArr.Databases); i++ {
		var database Database

		fmt.Println(dbDataArr.Databases[i].Type)
		switch dbType := dbDataArr.Databases[i].Type; dbType {
		case "mongo":
			database = &mongoDatabase{DB: &dbDataArr.Databases[i]}
		case "postgres":
			database = &postgresDatabase{DB: &dbDataArr.Databases[i]}
		default:
			database = nil
		}

		d.DBMap[dbDataArr.Databases[i].Name] = &database

		// fmt.Printf("val: %s\n", dbDataArr.Databases[i].Name)
	}
	// fmt.Println(d.DBMap)
	fmt.Println("----------------")
}

// ValidateYAML calls validation method on each database data object.
func (d *Databases) ValidateYAML() {
	fmt.Println("Database YAML file validation...")
	for _, database := range d.DBMap {
		(*database).GetData().Validate()
	}
	fmt.Println("...passed.")
}
