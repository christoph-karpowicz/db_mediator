package synch

type record struct {
	Data     map[string]interface{}
	Key      string
	ActiveIn []*vector
	PairedIn []*vector
}
