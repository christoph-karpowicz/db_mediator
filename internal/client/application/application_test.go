package application

import (
	"fmt"
	"os"
	"testing"
)

var app Application = Application{}

func printDebug(res map[string]interface{}) {
	for key, val := range res {
		switch val.(type) {
		case map[string]interface{}:
			fmt.Printf("----\n%s: \n", key)
			printDebug(val.(map[string]interface{}))
			fmt.Println("----")
		case []interface{}:
			fmt.Printf("----\n%s: \n", key)
			for subKey, subVal := range val.([]interface{}) {
				switch subVal.(type) {
				case map[string]interface{}:
					fmt.Printf("----\n%v: \n", subKey)
					printDebug(subVal.(map[string]interface{}))
					fmt.Println("----")
				default:
					fmt.Printf("	%s\n", subVal)
				}
			}
		case []map[string]interface{}:
			for _, subVal := range val.([]map[string]interface{}) {
				printDebug(subVal)
			}
		default:
			fmt.Printf("%s: %s\n", key, val)
		}
	}
}

func TestSimulationRequest(t *testing.T) {
	os.Chdir("../../..")
	app.runSynch("one-off", "films", true)
}
