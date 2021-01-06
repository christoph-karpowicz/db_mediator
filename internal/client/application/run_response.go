package application

import (
	"fmt"
)

// printRunResponse dispatches the response to the corresponding printer function.
func printRunResponse(res map[string]interface{}, resType string) {
	if res["err"].(bool) {
		runResponsePrinters["error"](res)
	} else {
		runResponsePrinters[resType](res)
	}
}

// runResponsePrinters is a map of functions that print the JSON responses received from the backend.
var runResponsePrinters map[string]func(map[string]interface{}) = map[string]func(map[string]interface{}){
	// Error printer.
	"error": func(res map[string]interface{}) {
		fmt.Println(res["message"].(string))
	},

	// one-off synch printer.
	"one-off": func(res map[string]interface{}) {
		fmt.Println(res["payload"].(string))
		fmt.Println(res["message"].(string))
	},

	// ongoing synch printer.
	"ongoing": func(res map[string]interface{}) {
		fmt.Println(res["payload"].(string))
		fmt.Println(res["message"].(string))
	},
}
