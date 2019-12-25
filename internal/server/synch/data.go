package synch

import (
	validationUtil "github.com/christoph-karpowicz/unifier/internal/server/util/validation"
)

type databases struct {
	Db1 database `json:"db1"`
	Db2 database `json:"db2"`
}

type database struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type synchData struct {
	Name      string    `json:"name"`
	Databases databases `json:"databases"`
	Tables    []table   `json:"tables"`
	Active    bool      `json:"active"`
}

func (s *synchData) Validate() {
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
			validationUtil.JSONStruct(vector.Conditions)
		}
	}
}
