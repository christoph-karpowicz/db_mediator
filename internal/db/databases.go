package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Databases struct {
	DBMap map[string]*Database
}

func (d *Databases) ImportJSON() {
	databasesFilePath, _ := filepath.Abs("./config/databases.json")

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

	var databases map[string]json.RawMessage
	var databasesArray []DatabaseData

	json.Unmarshal(byteArray, &databases)
	json.Unmarshal(databases["databases"], &databasesArray)

	fmt.Println("----------------")
	fmt.Println("Databases:")
	for i := 0; i < len(databasesArray); i++ {
		var database Database

		fmt.Println(databasesArray[i].Type)
		switch dbType := databasesArray[i].Type; dbType {
		case "mongo":
			database = &MongoDatabase{}
			dbb := database.(*MongoDatabase)
			dbb.Tst = 2222222222
			fmt.Println(dbb.Tst)
		case "postgres":
			database = &PostgresDatabase{}
		default:
			database = nil
		}

		d.DBMap[databasesArray[i].Name] = &database

		// fmt.Printf("key[%s] value[%s]\n", k, v)
		fmt.Println(database)
	}
	fmt.Println(d.DBMap)
	fmt.Println("----------------")
}
