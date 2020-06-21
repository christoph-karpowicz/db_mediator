package db

import (
	validationUtil "github.com/christoph-karpowicz/unifier/internal/util/validation"
)

var nullableFields = []string{"alias"}

// configArray is an array of YAML database configs.
type configArray struct {
	Databases []config
}

// config represents an individual YAML database config.
type config struct {
	Name     string `yaml:"name"`
	Alias    string `yaml:"alias"`
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// GetName returns the DB's name if an alias hasn't been provided.
func (d *config) GetName() string {
	if d.Alias != "" {
		return d.Alias
	}
	return d.Name
}

// Validate calls a validation function on itself.
func (d *config) Validate() {
	validationUtil.YAMLStruct(*d, nullableFields)
}
