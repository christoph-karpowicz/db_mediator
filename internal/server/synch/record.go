package synch

type record struct {
	Data     map[string]interface{}
	ActiveIn []*Mapping
	PairedIn []*Mapping
}
