package application

type response struct {
	Err bool   `json:"err"`
	Msg string `json:"msg"`
}

func createResponse(appRes interface{}) *response {
	var res *response

	switch appRes.(type) {
	case error:
		res = &response{
			Err: true,
			Msg: appRes.(error).Error(),
		}
	case []byte:
		res = &response{
			Err: false,
			Msg: string(appRes.([]byte)),
		}
	}

	return res
}
