package synch

type Flow struct {
	sourceColumnName string
	targetColumnName string
	source           *Record
	target           *Record
}
