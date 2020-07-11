package synch

type counters struct {
	fullSelects int
}

func newCounters() *counters {
	return &counters{0}
}
