package synch

import "encoding/json"

type operation interface {
	toJSON() string
}

type updateOrIdleOperation struct {
	IterationId       string      `json:"iterationId"`
	Timestamp         string      `json:"timestamp"`
	Operation         string      `json:"operation"`
	SourceTableName   string      `json:"sourceTableName"`
	SourceKeyName     string      `json:"sourceKeyName"`
	SourceKeyValue    interface{} `json:"sourceKeyValue"`
	SourceColumnName  string      `json:"sourceColumnName"`
	SourceColumnValue interface{} `json:"sourceColumnValue"`
	TargetKeyName     string      `json:"targetKeyName"`
	TargetKeyValue    interface{} `json:"targetKeyValue"`
	TargetColumnName  string      `json:"targetColumnName"`
	TargetColumnValue interface{} `json:"targetColumnValue"`
	TargetTableName   string      `json:"targetTableName"`
}

func (o *updateOrIdleOperation) toJSON() string {
	operationsJSON, err := json.MarshalIndent(o, "", "	")
	if err != nil {
		panic(err)
	}
	return string(operationsJSON)
}

type insertOperation struct {
	IterationId       string                 `json:"iterationId"`
	Timestamp         string                 `json:"timestamp"`
	Operation         string                 `json:"operation"`
	SourceTableName   string                 `json:"sourceTableName"`
	SourceKeyName     string                 `json:"sourceKeyName"`
	SourceKeyValue    interface{}            `json:"sourceKeyValue"`
	SourceColumnName  string                 `json:"sourceColumnName"`
	SourceColumnValue interface{}            `json:"sourceColumnValue"`
	TargetTableName   string                 `json:"targetTableName"`
	InsertedRow       map[string]interface{} `json:"insertedRow"`
}

func (o *insertOperation) toJSON() string {
	operationsJSON, err := json.MarshalIndent(o, "", "	")
	if err != nil {
		panic(err)
	}
	return string(operationsJSON)
}
