package synch

type Pair struct {
	record1    *Record
	record2    *Record
	IsComplete bool
}

func (p *Pair) Synchronize() {
	// return s.synch
}
