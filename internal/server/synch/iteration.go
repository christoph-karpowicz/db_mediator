package synch

import (
	"fmt"
	"strconv"
	"time"
)

type iteration struct {
	id         string
	synch      *Synch
	operations []operation
}

func newIteration(synch *Synch) *iteration {
	return &iteration{
		id:    getNewIterationID(synch),
		synch: synch,
	}
}

func getNewIterationID(synch *Synch) string {
	return synch.cfg.Name + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func (i *iteration) addOperation(op operation) {
	i.operations = append(i.operations, op)
	if !i.synch.IsSimulation() {
		fmt.Println(op.toJSON())
	}
}

func (i *iteration) flush() {
	i.synch.result.Operations = append(i.synch.result.Operations, i.operations...)
}
