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
	validationUtil.YAMLStruct(*s)
	validationUtil.YAMLStruct(s.Databases.Db1)
	validationUtil.YAMLStruct(s.Databases.Db2)
	for _, table := range s.Tables {
		validationUtil.YAMLStruct(table)

		validationUtil.YAMLStruct(table.Settings.SynchType)
		validationUtil.YAMLStruct(table.Settings.CreateNewRows)
		validationUtil.YAMLStruct(table.Settings.UpdateOldRows)

		for _, vector := range table.Vectors {
			validationUtil.YAMLStruct(vector)
			validationUtil.YAMLStruct(vector.ColumnNames)
			validationUtil.YAMLStruct(vector.Conditions)
		}
	}
}
