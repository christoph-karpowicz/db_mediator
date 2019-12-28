package db

import (
	validationUtil "github.com/christoph-karpowicz/unifier/internal/util/validation"
)

var nullableFields = []string{"alias"}

// DatabaseData reflects an array of YAML database configs.
type databaseDataArray struct {
	Databases []databaseData
}

// DatabaseData reflects an individual YAML database config.
type databaseData struct {
	Name     string `yaml:"name"`
	Alias    string `yaml:"alias"`
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Validate calls a validation function on itself.
func (d *databaseData) Validate() {
	validationUtil.YAMLStruct(*d, nullableFields)
}
