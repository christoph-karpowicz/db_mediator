package synch

import (
	validationUtil "github.com/christoph-karpowicz/unifier/internal/server/util/validation"
)

type Databases struct {
	Db1 Database `json:"db1"`
	Db2 Database `json:"db2"`
}

type Database struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SynchData struct {
	Name      string    `json:"name"`
	Databases Databases `json:"databases"`
	Tables    []Table   `json:"tables"`
	Active    bool      `json:"active"`
}

func (s *SynchData) Validate() {
	validationUtil.JSONStruct(*s)
	validationUtil.JSONStruct(s.Databases.Db1)
	validationUtil.JSONStruct(s.Databases.Db2)
	for _, table := range s.Tables {
		validationUtil.JSONStruct(table)

		validationUtil.JSONStruct(table.Settings.SynchType)
		validationUtil.JSONStruct(table.Settings.CreateNewRows)
		validationUtil.JSONStruct(table.Settings.UpdateOldRows)

		for _, vector := range table.Vectors {
			validationUtil.JSONStruct(vector)
			validationUtil.JSONStruct(vector.ColumnNames)
			validationUtil.JSONStruct(vector.Wheres)
		}
	}
}
