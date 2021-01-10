package synch

import (
	"fmt"
	"strconv"
	"time"
)

type iteration struct {
	id         string
	synch      *Synch
	operations []*operation
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

func (i *iteration) addOperation(operation *operation) {
	i.operations = append(i.operations, operation)
	if !i.synch.IsSimulation() {
		fmt.Println(operation.toJSON())
	}
}

func (i *iteration) flush() {
	i.synch.result.Operations = append(i.synch.result.Operations, i.operations...)
}
