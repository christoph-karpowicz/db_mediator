package application

import "fmt"

// printStopResponse dispatches the response to the corresponding printer function.
func printStopResponse(res map[string]interface{}) {
	// var resStr string

	// if res["err"].(bool) {
	// 	resStr = responsePrinters["error"](res)
	// } else {
	// 	resStr = responsePrinters[resType](res)
	// }

	// fmt.Println(resStr)
	fmt.Println(res)
}
