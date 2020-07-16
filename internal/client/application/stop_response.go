package application

import "fmt"

// printStopResponse dispatches the response to the corresponding printer function.
func printStopResponse(res map[string]interface{}) {
	var resStr string

	if res["err"].(bool) {
		resStr = stopResponsePrinters["error"](res)
	} else {
		resStr = stopResponsePrinters["stop"](res)
	}

	fmt.Println(resStr)
}

// stopResponsePrinters is a map of functions that print the JSON responses received from the backend.
var stopResponsePrinters map[string]func(map[string]interface{}) string = map[string]func(map[string]interface{}) string{

	// Error printer.
	"error": func(res map[string]interface{}) string {
		return res["payload"].(string)
	},

	// ongoing synch printer.
	"stop": func(res map[string]interface{}) string {
		return res["payload"].(string)
	},
}
