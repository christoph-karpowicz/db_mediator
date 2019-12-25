package synch

type flow struct {
	sourceColumnName string
	targetColumnName string
	source           *record
	target           *record
}
