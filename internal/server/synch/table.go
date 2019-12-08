package synch

type TableSpecifics struct {
	Table1 string `json:"table1"`
	Table2 string `json:"table2"`
}

type SynchType struct {
	MatchBy     string         `json:"matchBy"`
	ColumnNames TableSpecifics `json:"columnNames"`
}

type Settings struct {
	SynchType     SynchType
	CreateNewRows TableSpecifics
	UpdateOldRows TableSpecifics
}

type Table struct {
	Names    TableSpecifics `json:"names"`
	Keys     TableSpecifics `json:"primaryKeys"`
	Settings Settings       `json:"settings"`
	Vectors  []Vector       `json:"vectors"`
	Db1Data  []map[string]interface{}
	Db2Data  []map[string]interface{}
}
