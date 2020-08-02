package application

import (
	"encoding/json"
)

const (
	PAYLOAD_TYPE_JSON = "json"
	PAYLOAD_TYPE_TEXT = "text"
)

// parseResponse turns a JSON backend reponse to a processable map.
func parseResponse(res []byte) map[string]interface{} {
	result := make(map[string]interface{})

	if err := json.Unmarshal(res, &result); err != nil {
		panic(err)
	}

	return result
}
