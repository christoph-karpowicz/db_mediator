package synch

import "fmt"

type synchInitError struct {
	method string
	errMsg string
}

func (e *synchInitError) Error() string {
	return fmt.Sprintf("[synch init] %s in method %s", e.errMsg, e.method)
}
