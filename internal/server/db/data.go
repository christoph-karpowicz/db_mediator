package db

import (
	validationUtil "github.com/christoph-karpowicz/unifier/internal/server/util/validation"
)

// DatabaseData reflects JSON database config.
type databaseData struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Validate calls a validation function on itself.
func (d *databaseData) Validate() {
	validationUtil.JSONStruct(*d)
}
