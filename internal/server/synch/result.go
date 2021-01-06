package synch

import (
	"encoding/json"
)

type Result struct {
	Message    string       `json:"message"`
	Operations []*operation `json:"operations"`
}

func (r *Result) OperationsToJSON() string {
	operationsJSON, err := json.MarshalIndent(r.Operations, "", "	")
	if err != nil {
		panic(err)
	}
	return string(operationsJSON)
}
