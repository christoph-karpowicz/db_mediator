package synch

type flow struct {
	sourceTableName  string
	targetTableName  string
	sourceColumnName string
	targetColumnName string
	source           *record
	target           *record
}
