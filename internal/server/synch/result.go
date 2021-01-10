package synch

import (
	"encoding/json"
)

const (
	SIMULATION_DIR = "./simulation/"
	LOGS_DIR       = "./log/"
)

type Result struct {
	Message    string       `json:"message"`
	Operations []*operation `json:"operations"`
	path       string
}

func (r *Result) OperationsToJSON() string {
	operationsJSON, err := json.MarshalIndent(r.Operations, "", "	")
	if err != nil {
		panic(err)
	}
	return string(operationsJSON)
}

func (r *Result) operationsToJSONSlice() []string {
	operationsToJSON := make([]string, 0)
	for _, operation := range r.Operations {
		operationJSON := operation.toJSON()
		// return false, &SynchReportError{SynchName: r.synch.GetConfig().Name, ErrMsg: err.Error()}
		operationsToJSON = append(operationsToJSON, string(operationJSON))
	}
	return operationsToJSON
}

func (r *Result) setSimulationPath(fileName string) {
	r.path = SIMULATION_DIR + fileName
}

func (r *Result) setLogPath(fileName string) {
	r.path = LOGS_DIR + fileName
}
