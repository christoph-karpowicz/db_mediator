package application

import "fmt"

// printStopResponse dispatches the response to the corresponding printer function.
func printStopResponse(res map[string]interface{}) {
	var resStr string

	fmt.Println(res["payloadType"].(string))
	if res["err"].(bool) {
		resStr = stopResponsePrinters["error"](res)
	} else {
		if res["payloadType"].(string) == PAYLOAD_TYPE_TEXT {
			resStr = stopResponsePrinters["stop"](res)
		} else if res["payloadType"].(string) == PAYLOAD_TYPE_JSON {
			resStr = runResponsePrinters["one-off"](res)
		}
	}

	fmt.Println(resStr)
	fmt.Println(res["payload"].(string))
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
