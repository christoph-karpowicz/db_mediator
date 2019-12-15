package synch

type Vector struct {
	ColumnNames      TableSpecifics `json:"columnNames"`
	DataFlow         string         `json:"dataFlow"`
	Conditions       TableSpecifics `json:"conditions"`
	Db1ActiveRecords []*Record
	Db2ActiveRecords []*Record
	Pairs            []Pair
}
