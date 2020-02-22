package application

import (
	"encoding/json"
	"fmt"
	"strconv"
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
	resStr := responsePrinters[resType](res)
	fmt.Println(resStr)
}

// responsePrinters is a map of functions that print the JSON responses received from the backend.
var responsePrinters map[string]func(map[string]interface{}) string = map[string]func(map[string]interface{}) string{

	// Simulation printer.
	"simulation": func(res map[string]interface{}) string {
		synchInfo := res["synchInfo"].(map[string]interface{})
		mappingsStrArray := synchInfo["Mappings"].([]interface{})
		mappings := res["mappings"].(map[string]interface{})

		var allMappingsStr string
		for key, mapping := range mappings {
			mappingLinks := mapping.(map[string]interface{})["links"].(map[string]interface{})
			keyInt, err := strconv.Atoi(key)
			if err != nil {
				panic(err)
			}

			mappingStr := fmt.Sprintf(`
	Mapping index: %s
	command: %s
	links:
			`,
				key,
				mappingsStrArray[keyInt],
			)

			for linkKey, link := range mappingLinks {
				linkStr := fmt.Sprintf(`
		Link index: %s
`,
					linkKey,
				)

				for action, actionStrings := range link.(map[string]interface{}) {
					linkStr += fmt.Sprintf(`			%s:
`,
						action,
					)

					if actionStrings != nil {
						for _, actionString := range actionStrings.([]interface{}) {
							linkStr += fmt.Sprintf(`%s`,
								actionString.(string),
							)

						}
					}

				}

				mappingStr += linkStr
			}

			allMappingsStr += mappingStr
		}

		return fmt.Sprintf(`Synch name: %s
Mappings:
%s
		`,
			synchInfo["Name"].(string),
			allMappingsStr,
		)
	},
}
