package synch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	SIMULATION_DIR = "./simulation/"
	LOGS_DIR       = "./log/"
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
}

func (i *iteration) flush() *Result {
	iterationResult := Result{}
	operationsToJSON, err := i.operationsToJSON()
	if err != nil {
		panic(err)
	}
	operationsToJSONString := strings.Join(operationsToJSON, "\n")
	if i.synch.IsSimulation() {
		fileName := SIMULATION_DIR + i.id
		err := ioutil.WriteFile(fileName, []byte(operationsToJSONString), 0644)
		if err != nil {
			panic(err)
		}
		iterationResult.Message = fmt.Sprintf("Simulation report saved to file %s", fileName)
	} else {
		fmt.Println(operationsToJSONString)
	}
	iterationResult.Operations = i.operations
	return &iterationResult
}

func (i *iteration) operationsToJSON() ([]string, error) {
	operationsToJSON := make([]string, 0)
	for _, operation := range i.operations {
		operationJSON, err := json.MarshalIndent(&operation, "", "	")
		if err != nil {
			return nil, err
			// return false, &SynchReportError{SynchName: r.synch.GetConfig().Name, ErrMsg: err.Error()}
		}
		operationsToJSON = append(operationsToJSON, string(operationJSON))
	}
	return operationsToJSON, nil
}
