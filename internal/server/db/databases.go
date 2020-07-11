/*
Package db contains database configurations and
methods for querying.
*/
package db

import (
	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
)

// Databases imports, validates and holds information about databases from JSON config files.
type Databases map[string]*Database

func (d *Databases) Init() {
	dbCfgs := d.getConfigs()
	d.validateConfigs(dbCfgs)
}

func (d *Databases) getConfigs() *cfg.DbConfigArray {
	var dbCfgs *cfg.DbConfigArray = cfg.GetDbConfigs()

	for i := 0; i < len(dbCfgs.Databases); i++ {
		var database Database

		// fmt.Println(dbCfgs.Databases[i].Type)
		switch dbType := dbCfgs.Databases[i].Type; dbType {
		case "mongo":
			database = &mongoDatabase{cfg: &dbCfgs.Databases[i]}
		case "postgres":
			database = &postgresDatabase{cfg: &dbCfgs.Databases[i]}
		default:
			database = nil
		}

		(*d)[dbCfgs.Databases[i].GetName()] = &database
		// fmt.Printf("val: %s\n", dbDataArr.Databases[i].Name)
	}

	return dbCfgs
}

// validateConfigs calls validation method on each database data object.
func (d *Databases) validateConfigs(dbCfgs *cfg.DbConfigArray) {
	dbCfgs.Validate()
}
