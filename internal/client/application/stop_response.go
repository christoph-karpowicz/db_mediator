package application

import "fmt"

// printStopResponse dispatches the response to the corresponding printer function.
func printStopResponse(res map[string]interface{}) {

	fmt.Println(res["payloadType"].(string))
	if res["err"].(bool) {
		stopResponsePrinters["error"](res)
	} else {
		if res["payloadType"].(string) == PAYLOAD_TYPE_TEXT {
			stopResponsePrinters["stop"](res)
		} else if res["payloadType"].(string) == PAYLOAD_TYPE_JSON {
			runResponsePrinters["one-off"](res)
		}
	}
}

// stopResponsePrinters is a map of functions that print the JSON responses received from the backend.
var stopResponsePrinters map[string]func(map[string]interface{}) = map[string]func(map[string]interface{}){

	// Error printer.
	"error": func(res map[string]interface{}) {
		fmt.Println(res["payload"].(string))
	},

	// ongoing synch printer.
	"stop": func(res map[string]interface{}) {
		fmt.Println(res["payload"].(string))
	},
}
