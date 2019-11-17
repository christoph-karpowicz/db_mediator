package synch

import (
	validationUtil "github.com/christoph-karpowicz/unifier/internal/util/validation"
)

type SynchData struct {
	Name      string   `json:"name"`
	Databases []string `json:"databases"`
	Tables    []Table  `json:"tables"`
}

func (s *SynchData) Validate() {
	validationUtil.JSONStruct(*s, "")
	for _, table := range s.Tables {
		validationUtil.JSONStruct(table, "")
		var ignoreInConnection string
		switch table.Connection.Type {
		case "external_id_columns":
			ignoreInConnection = "InitialVector"
		case "persistence":
			ignoreInConnection = "Columns"
		default:
			ignoreInConnection = ""
		}
		validationUtil.JSONStruct(table.Connection, ignoreInConnection)
	}
}
