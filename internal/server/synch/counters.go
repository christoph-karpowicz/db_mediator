package synch

type counters struct {
	selects int
}

func newCounters() *counters {
	return &counters{0}
}

func (c *counters) reset() {
	c.selects = 0
}
