package application

import (
	"encoding/json"
	"fmt"
	"strings"
)

// parseResponse turns a JSON backend reponse to a processable map.
func parseResponse(res []byte) map[string]interface{} {
	result := make(map[string]interface{})

	if err := json.Unmarshal(res, &result); err != nil {
		panic(err)
	}
	// log.Println(result)

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

// printAction converts a JSON action representation into a readable string.
func printAction(actionString string) string {
	actionMap := make(map[string]interface{})

	if err := json.Unmarshal([]byte(actionString), &actionMap); err != nil {
		panic(err)
	}
	// log.Println(actionMap)

	return actionPrinters[actionMap["actionType"].(string)](actionMap)
}

var actionPrinters map[string]func(map[string]interface{}) string = map[string]func(map[string]interface{}) string{
	// Creates a string representation of two records that
	// are the same and no action will be carried out.
	"idle": func(action map[string]interface{}) string {
		return fmt.Sprintf("|%6v: %3v, %6v: %25v|  ==  |%6v: %6v, %6s: %25v|\n",
			action["sourceNodeKey"],
			action["sourceData"],
			action["sourceColumn"],
			action["sourceColumnData"],
			action["targetKeyName"],
			action["targetKeyValue"],
			action["targetColumn"],
			action["targetColumnData"],
		)
	},
	// Creates a string representation of an insert
	// that would be carried out due to the pair's incompleteness.
	"insert": func(action map[string]interface{}) string {
		return fmt.Sprintf("|%6v: %3v, %6v: %25v|  =>  |%6v: %6v, %6s: %25v|\n",
			action["sourceNodeKey"],
			action["sourceData"],
			action["sourceColumn"],
			action["sourceColumnData"],
			action["targetKeyName"],
			"-",
			action["targetColumn"],
			action["sourceColumnData"],
		)
	},
	// Creates a string representation of an update
	// that would be carried out because the data in the pair's records
	// was found to be different.
	"update": func(action map[string]interface{}) string {
		return fmt.Sprintf("|%6v: %3v, %6v: %25v|  =^  |%6v: %6v, %6s: %25v -> %25v|\n",
			action["sourceNodeKey"],
			action["sourceData"],
			action["sourceColumn"],
			action["sourceColumnData"],
			action["targetKeyName"],
			action["targetKeyValue"],
			action["targetColumn"],
			action["targetColumnData"],
			action["sourceColumnData"],
		)
	},
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
		links := resPayload["links"].(map[string]interface{})

		// LINKS
		linksStr := ""
		for linkIndex, link := range links {
			linkStr := fmt.Sprintf(`
=================
Link index: %s
Link command: %s`,
				linkIndex,
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
							printAction(actionString.(string)),
						)

					}
				}

				linkStr += fmt.Sprintf("%s",
					horizontalBorder,
				)

			}

			linksStr += linkStr
		}

		// WHOLE SIMULATION
		return fmt.Sprintf(`SYNCH NAME: %s
SERVER RESPONSE: %s

LINKS:%s
		`,
			synchInfo["Name"].(string),
			synchMsg,
			linksStr,
		)
	},
}
