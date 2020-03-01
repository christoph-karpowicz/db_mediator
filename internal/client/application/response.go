package application

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// parseResponse turns a JSON backend reponse to a processable map.
func parseResponse(res []byte) map[string]interface{} {
	result := make(map[string]interface{})

	if err := json.Unmarshal(res, &result); err != nil {
		panic(err)
	}

	return result
}

// printResponse dispatches the response to the corresponding printer function.
func printResponse(res map[string]interface{}, resType string) {
	var resStr string

	if res["err"].(bool) {
		resStr = responsePrinters["error"](res)
	} else {
		resStr = responsePrinters[resType](res)
	}

	fmt.Println(resStr)
}

// responsePrinters is a map of functions that print the JSON responses received from the backend.
var responsePrinters map[string]func(map[string]interface{}) string = map[string]func(map[string]interface{}) string{

	// Error printer.
	"error": func(res map[string]interface{}) string {
		return res["payload"].(string)
	},

	// Simulation printer.
	"one-off": func(res map[string]interface{}) string {
		resPayloadStr := res["payload"].(string)
		resPayload := make(map[string]interface{})

		if err := json.Unmarshal([]byte(resPayloadStr), &resPayload); err != nil {
			panic(err)
		}

		synchMsg := resPayload["msg"].(string)
		synchInfo := resPayload["synchInfo"].(map[string]interface{})
		mappingsStrArray := synchInfo["Mappings"].([]interface{})
		mappings := resPayload["mappings"].(map[string]interface{})

		// MAPPINGS
		var allMappingsStr string
		for key, mapping := range mappings {
			mappingLinks := mapping.(map[string]interface{})["links"].(map[string]interface{})
			keyInt, err := strconv.Atoi(key)
			if err != nil {
				panic(err)
			}

			mappingStr := fmt.Sprintf(`
==================================
Mapping index: %s
command: %s
links:
			`,
				key,
				mappingsStrArray[keyInt],
			)

			// LINKS
			for linkKey, link := range mappingLinks {
				linkStr := fmt.Sprintf(`
=================
Link index: %s
Link command: %s`,
					linkKey,
					link.(map[string]interface{})["cmd"].(string),
				)

				// ACTIONS
				for action, actionStrings := range link.(map[string]interface{}) {
					if action == "cmd" {
						continue
					}

					var horizontalBorder string
					switch action {
					case "updates":
						horizontalBorder = strings.Repeat("-", 64*2+7)
					default:
						horizontalBorder = strings.Repeat("-", 50*2+6)
					}

					linkStr += fmt.Sprintf(`
%s:`,
						strings.ToUpper(action),
					)

					if actionStrings == nil {
						linkStr += fmt.Sprintf("\nno %s actions", action)
						continue
					} else {
						linkStr += fmt.Sprintf("\n%s\n", horizontalBorder)
					}

					// ACTUAL MODIFICATIONS OF RECORDS
					if actionStrings != nil {
						for _, actionString := range actionStrings.([]interface{}) {
							linkStr += fmt.Sprintf("%s",
								actionString.(string),
							)

						}
					}

					linkStr += fmt.Sprintf("%s",
						horizontalBorder,
					)

				}

				mappingStr += linkStr
			}

			allMappingsStr += mappingStr
		}

		// WHOLE SIMULATION
		return fmt.Sprintf(`SYNCH NAME: %s
SERVER RESPONSE: %s
MAPPINGS:
%s
		`,
			synchInfo["Name"].(string),
			synchMsg,
			allMappingsStr,
		)
	},
}
