package synch

type Record struct {
	Data     map[string]interface{}
	ActiveIn []*Vector
	PairedIn []*Vector
}
