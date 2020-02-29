package application

type response struct {
	Success bool   `json:"success"`
	Err     bool   `json:"err"`
	Msg     string `json:"msg"`
}

func createResponse(appRes interface{}) *response {
	var res *response

	switch appRes.(type) {
	case error:
		res = &response{
			Success: false,
			Err:     true,
			Msg:     appRes.(error).Error(),
		}
	case []byte:
		res = &response{
			Success: true,
			Err:     false,
			Msg:     string(appRes.([]byte)),
		}
	}

	return res
}
