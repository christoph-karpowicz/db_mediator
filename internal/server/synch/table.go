package synch

type tableSpecifics struct {
	Table1 string `json:"table1"`
	Table2 string `json:"table2"`
}

type synchType struct {
	MatchBy     string         `json:"matchBy"`
	ColumnNames tableSpecifics `json:"columnNames"`
}

type settings struct {
	SynchType     synchType
	CreateNewRows tableSpecifics
	UpdateOldRows tableSpecifics
}

type table struct {
	Names         tableSpecifics `json:"names"`
	Keys          tableSpecifics `json:"primaryKeys"`
	Settings      settings       `json:"settings"`
	Vectors       []vector       `json:"vectors"`
	Db1OldRecords *tableRecords
	Db2OldRecords *tableRecords
	Db1Records    *tableRecords
	Db2Records    *tableRecords
}
