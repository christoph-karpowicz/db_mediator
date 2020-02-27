package lang

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ParseMapping uses regexp to split the mapping string into smaller parts.
func ParseMapping(str string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	mappingParseRegexp := regexp.MustCompile(`(?isU)^\s*(?P<command>[a-z]+)\s+(?P<links>\[.*\]\sTO\s\[.*\],?)+\s+MATCH\sBY\s(?P<matchMethod>[a-z]+\(.+\))\s+DO\s(?P<do>[a-z\s,]+)\s*$`)
	matches := mappingParseRegexp.FindStringSubmatch(str)
	subNames := mappingParseRegexp.SubexpNames()

	commaSepRegexp := regexp.MustCompile(`(?s)\s*,\s*`)

	for i, match := range matches {
		// Skip the first, empty element.
		if i == 0 {
			continue
		}

		switch subNames[i] {
		case "links":
			result[subNames[i]] = make([]map[string]string, 0)
			for _, link := range commaSepRegexp.Split(match, -1) {
				result[subNames[i]] = append(result[subNames[i]].([]map[string]string), parseLink(link))
			}
		case "matchMethod":
			result[subNames[i]] = parseMatchMethod(match)
		case "do":
			result[subNames[i]] = make([]string, 0)
			for _, doCmd := range commaSepRegexp.Split(match, -1) {
				result[subNames[i]] = append(result[subNames[i]].([]string), doCmd)
			}
			result[subNames[i]] = commaSepRegexp.Split(match, -1)
		default:
			result[subNames[i]] = match
		}
	}

	err := validateMapping(result)

	return result, err
}

// parseLink splits an individual link into smaller parts.
func parseMatchMethod(str string) map[string]interface{} {
	parsedMatchMethod := make(map[string]interface{})
	r := regexp.MustCompile(`(?iU)^(?P<matchCmd>[a-z]+)\((?P<matchArgs>.+)\)$`)
	matches := r.FindStringSubmatch(str)
	subNames := r.SubexpNames()

	commaSepRegexp := regexp.MustCompile(`(?s)\s*,\s*`)
	dotSepRegexp := regexp.MustCompile(`(?s)\.`)

	for i, match := range matches {
		if subNames[i] == "matchArgs" {
			parsedMatchMethod[subNames[i]] = make([]string, 0)
			for _, matchArg := range commaSepRegexp.Split(match, -1) {
				parsedMatchMethod[subNames[i]] = append(parsedMatchMethod[subNames[i]].([]string), matchArg)
			}
		} else if subNames[i] != "" {
			parsedMatchMethod[subNames[i]] = match
		}
	}

	// Extract the node names and external ID column names from match method.
	if parsedMatchMethod["matchCmd"].(string) == "IDS" {
		parsedMatchMethod["parsedMatchArgs"] = make([]map[string]string, 0)

		for _, matchArg := range parsedMatchMethod["matchArgs"].([]string) {
			parsedArg := make(map[string]string)
			splitArg := dotSepRegexp.Split(matchArg, -1)
			parsedArg["node"] = splitArg[0]
			parsedArg["extIDColumn"] = splitArg[1]
			parsedMatchMethod["parsedMatchArgs"] = append(parsedMatchMethod["parsedMatchArgs"].([]map[string]string), parsedArg)
		}
	}

	return parsedMatchMethod
}

// parseLink splits an individual link into smaller parts.
func parseLink(str string) map[string]string {
	parsedLink := make(map[string]string)
	r := regexp.MustCompile(`(?iU)^\[(?P<sourceNode>.+)\.(?P<sourceColumn>.+)(\s+WHERE\s+(?P<sourceWhere>.+))?\]\sTO\s\[(?P<targetNode>.+)\.(?P<targetColumn>.+)(\s+WHERE\s+(?P<targetWhere>.+))?\]$`)
	matches := r.FindStringSubmatch(str)
	subNames := r.SubexpNames()

	parsedLink["raw"] = str
	for i, match := range matches {
		if subNames[i] != "" {
			parsedLink[subNames[i]] = match
		}
	}

	return parsedLink
}

func validateMapping(result map[string]interface{}) error {
	errorsArr := make([]string, 0)
	var err error = nil
	var preError string = "[mapping parsing] ERROR: "

	fmt.Println(result)

	// command
	if result["command"] == nil || len(result["command"].(string)) == 0 {
		errorsArr = append(errorsArr, preError+"no command found")
	}

	// links
	if result["links"] == nil || len(result["links"].([]map[string]string)) == 0 {
		errorsArr = append(errorsArr, preError+"no links found")
	}

	if len(errorsArr) > 0 {
		errorsArrJoined := strings.Join(errorsArr, "\n")
		err = errors.New(errorsArrJoined)
	}
	return err
}
