package synch

type Record struct {
	Data     map[string]interface{}
	Key      string
	ActiveIn []*Vector
	PairedIn []*Vector
}
