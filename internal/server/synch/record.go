package synch

type record struct {
	Data     map[string]interface{}
	ActiveIn []*Link
	PairedIn []*Link
}
