package synch

type action struct {
	linkId           string
	atype            string
	SourceNodeKey    string      `json:"sourceNodeKey"`
	SourceData       interface{} `json:"sourceData"`
	SourceColumn     string      `json:"sourceColumn"`
	SourceColumnData interface{} `json:"sourceColumnData"`
	TargetKeyName    string      `json:"targetKeyName"`
	TargetKeyValue   interface{} `json:"targetKeyValue"`
	TargetColumn     string      `json:"targetColumn"`
	TargetColumnData interface{} `json:"targetColumnData"`
	ActionType       string      `json:"actionType"`
}
