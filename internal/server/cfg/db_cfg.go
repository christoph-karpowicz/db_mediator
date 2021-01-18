package cfg

import (
	validationUtil "github.com/christoph-karpowicz/db_mediator/internal/util/validation"
)

var dbNullableFields = []string{"alias"}

// DbConfigArray is an array of YAML database configs.
type DbConfigArray struct {
	Databases []DbConfig
}

// DbConfig represents an individual YAML database config.
type DbConfig struct {
	Name     string `yaml:"name"`
	Alias    string `yaml:"alias"`
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// GetName returns the DB's name if an alias hasn't been provided.
func (d *DbConfig) GetName() string {
	if d.Alias != "" {
		return d.Alias
	}
	return d.Name
}

// Validate calls a validation function on itself.
func (d *DbConfigArray) Validate() {
	for _, dbCfg := range d.Databases {
		validationUtil.YAMLStruct(dbCfg, dbNullableFields)
	}
}

// GetDbConfigs loads the database configs from databases.yaml file.
func GetDbConfigs() *DbConfigArray {
	var dbDataArr DbConfigArray = DbConfigArray{}
	ImportYAMLFile(&dbDataArr, "./config/databases.yaml")
	return &dbDataArr
}
