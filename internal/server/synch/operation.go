package synch

import "encoding/json"

type operation struct {
	IterationId       string      `json:"iterationId"`
	Timestamp         string      `json:"timestamp"`
	Operation         string      `json:"operation"`
	SourceKeyName     string      `json:"sourceKeyName"`
	SourceKeyValue    interface{} `json:"sourceKeyValue"`
	SourceColumnName  string      `json:"sourceColumnName"`
	SourceColumnValue interface{} `json:"sourceColumnValue"`
	TargetKeyName     string      `json:"targetKeyName"`
	TargetKeyValue    interface{} `json:"targetKeyValue"`
	TargetColumnName  string      `json:"targetColumnName"`
	TargetColumnValue interface{} `json:"targetColumnValue"`
}

func (o *operation) toJSON() string {
	operationsJSON, err := json.MarshalIndent(o, "", "	")
	if err != nil {
		panic(err)
	}
	return string(operationsJSON)
}
