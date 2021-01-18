package application

import "github.com/christoph-karpowicz/db_mediator/internal/server/synch"

type response struct {
	Err     bool   `json:"err"`
	Message string `json:"message"`
	Payload string `json:"payload"`
}

func createResponseChannel() chan *response {
	return make(chan *response)
}

func createResponse(synchResult interface{}) *response {
	var res *response

	switch synchResult.(type) {
	case error:
		panic(synchResult.(error))
		res = &response{
			Err:     true,
			Message: synchResult.(error).Error(),
		}
	case string:
		res = &response{
			Err:     false,
			Message: synchResult.(string),
		}
	case *synch.Result:
		res = &response{
			Err:     false,
			Message: synchResult.(*synch.Result).Message,
			Payload: synchResult.(*synch.Result).OperationsToJSON(),
		}
	}

	return res
}
