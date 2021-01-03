package application

import (
	"fmt"
)

// printRunResponse dispatches the response to the corresponding printer function.
func printRunResponse(res map[string]interface{}, resType string) {
	var resStr string

	if res["err"].(bool) {
		resStr = runResponsePrinters["error"](res)
	} else {
		resStr = runResponsePrinters[resType](res)
	}

	fmt.Println(resStr)
}

// runResponsePrinters is a map of functions that print the JSON responses received from the backend.
var runResponsePrinters map[string]func(map[string]interface{}) string = map[string]func(map[string]interface{}) string{
	// Error printer.
	"error": func(res map[string]interface{}) string {
		return res["payload"].(string)
	},

	// one-off synch printer.
	"one-off": func(res map[string]interface{}) string {
		return res["payload"].(string)
	},

	// ongoing synch printer.
	"ongoing": func(res map[string]interface{}) string {
		return res["payload"].(string)
	},
}
