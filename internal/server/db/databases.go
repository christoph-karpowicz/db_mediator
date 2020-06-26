/*
Package db contains database configurations and
methods for querying.
*/
package db

import (
	"fmt"

	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
)

// Databases imports, validates and holds information about databases from JSON config files.
type Databases struct {
	dbCfgs *cfg.DbConfigArray
	DBMap  map[string]*Database
}

func (d *Databases) Init() {
	d.getConfigs()
	d.validateConfigs()
	d.assignConfigs()
}

func (d *Databases) getConfigs() {
	var dbConfigArray *cfg.DbConfigArray = cfg.GetDbConfig()
	d.dbCfgs = dbConfigArray
}

func (d *Databases) assignConfigs() {
	for i := 0; i < len(d.dbCfgs.Databases); i++ {
		var database Database

		fmt.Println(d.dbCfgs.Databases[i].Type)
		switch dbType := d.dbCfgs.Databases[i].Type; dbType {
		case "mongo":
			database = &mongoDatabase{cfg: &d.dbCfgs.Databases[i]}
		case "postgres":
			database = &postgresDatabase{cfg: &d.dbCfgs.Databases[i]}
		default:
			database = nil
		}

		d.DBMap[d.dbCfgs.Databases[i].GetName()] = &database

		// fmt.Printf("val: %s\n", dbDataArr.Databases[i].Name)
	}
}

// validateConfigs calls validation method on each database data object.
func (d *Databases) validateConfigs() {
	fmt.Println("Database YAML file validation...")
	d.dbCfgs.Validate()
	fmt.Println("...passed.")
}
