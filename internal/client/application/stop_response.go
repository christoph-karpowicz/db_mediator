package application

import "fmt"

// printStopResponse dispatches the response to the corresponding printer function.
func printStopResponse(res map[string]interface{}) {
	fmt.Println(res["payload"].(string))
	fmt.Println(res["message"].(string))
}
