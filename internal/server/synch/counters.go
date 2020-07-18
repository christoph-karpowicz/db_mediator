package synch

type counters struct {
	fullSelects int
}

func newCounters() *counters {
	return &counters{0}
}

func (c *counters) reset() {
	c.fullSelects = 0
}
