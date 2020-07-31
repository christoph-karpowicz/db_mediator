package application

const (
	PAYLOAD_TYPE_JSON = "json"
	PAYLOAD_TYPE_TEXT = "text"
)

type response struct {
	Err         bool   `json:"err"`
	PayloadType string `json:"payloadType"`
	Payload     string `json:"payload"`
}

func createResponse(appRes interface{}) *response {
	var res *response

	switch appRes.(type) {
	case error:
		panic(appRes.(error))
		res = &response{
			Err:         true,
			PayloadType: PAYLOAD_TYPE_TEXT,
			Payload:     appRes.(error).Error(),
		}
	case string:
		res = &response{
			Err:         false,
			PayloadType: PAYLOAD_TYPE_TEXT,
			Payload:     appRes.(string),
		}
	case []byte:
		res = &response{
			Err:         false,
			PayloadType: PAYLOAD_TYPE_JSON,
			Payload:     string(appRes.([]byte)),
		}
	}

	return res
}
