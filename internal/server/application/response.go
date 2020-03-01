package application

type response struct {
	Err     bool   `json:"err"`
	Payload string `json:"payload"`
}

func createResponse(appRes interface{}) *response {
	var res *response

	switch appRes.(type) {
	case error:
		res = &response{
			Err:     true,
			Payload: appRes.(error).Error(),
		}
	case []byte:
		res = &response{
			Err:     false,
			Payload: string(appRes.([]byte)),
		}
	}

	return res
}
