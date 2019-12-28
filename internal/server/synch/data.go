package synch

import (
	validationUtil "github.com/christoph-karpowicz/unifier/internal/util/validation"
)

var nullableFields = []string{"alias"}

var synchConnectionTypes = [2]string{"external_id_columns", "persistence"}
var createNewRows = [3]string{"never", "initially", "always"}
var updateOldRows = [3]string{"never", "initially", "always"}

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
	validationUtil.YAMLStruct(*s, nullableFields)
	validationUtil.YAMLStruct(s.Databases.Db1, nullableFields)
	validationUtil.YAMLStruct(s.Databases.Db2, nullableFields)
	for _, table := range s.Tables {
		validationUtil.YAMLStruct(table, nullableFields)

		validationUtil.YAMLStruct(table.Settings.SynchType, nullableFields)
		validationUtil.YAMLStruct(table.Settings.CreateNewRows, nullableFields)
		validationUtil.YAMLStruct(table.Settings.UpdateOldRows, nullableFields)

		for _, vector := range table.Vectors {
			validationUtil.YAMLStruct(vector, nullableFields)
			validationUtil.YAMLStruct(vector.ColumnNames, nullableFields)
			validationUtil.YAMLStruct(vector.Conditions, nullableFields)
		}
	}
}
