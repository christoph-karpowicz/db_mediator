package synch

type record struct {
	Data     map[string]interface{}
	ActiveIn []*mapping
	PairedIn []*mapping
}
