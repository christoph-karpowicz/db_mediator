package db

import (
	validationUtil "github.com/christoph-karpowicz/unifier/internal/server/util/validation"
)

type DatabaseData struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (d *DatabaseData) Validate() {
	validationUtil.JSONStruct(*d)
}
