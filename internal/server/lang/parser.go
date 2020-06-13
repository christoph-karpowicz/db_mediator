package lang

import (
	"fmt"
	"regexp"
	"strings"

	arrUtil "github.com/christoph-karpowicz/unifier/internal/util/array"
)

/*
TODO:
1. change mapping parsing so if an error occurs,
	it could specify where exactly.
*/

type mappingParserError struct {
	errMsg string
}

func (e *mappingParserError) Error() string {
	return fmt.Sprintf("[mapping parser] %s", e.errMsg)
}

// ParseLink uses regexp to split the link string into smaller parts.
func ParseLink(str string) (map[string]string, error) {
	result := make(map[string]string)
	regexpString := `(?ismU)^\s*\[(?P<sourceNode>[^\.,]+)\.(?P<sourceColumn>[^\.,]+)(?P<sourceWhere>\s+WHERE\s+.+)?\]\s+TO\s+\[(?P<targetNode>[^\.,]+)\.(?P<targetColumn>[^\.,]+)(?P<targetWhere>\s+WHERE\s+.+)?\]\s*$`
	parseRegexp := regexp.MustCompile(regexpString)
	matches := parseRegexp.FindStringSubmatch(str)
	subNames := parseRegexp.SubexpNames()

	fmt.Println(matches)

	for i, match := range matches {
		// Skip the first, empty element.
		if i == 0 {
			continue
		}

		fmt.Println(match)

		if arrUtil.Contains([]string{"sourceWhere", "targetWhere"}, subNames[i]) {
			parsedWhere := ParseLinkWhere(match)
			result[subNames[i]] = parsedWhere
		} else {
			result[subNames[i]] = match
		}
	}

	// err := validateMapping(result)

	// return result, err
	return result, nil
}

// ParseLinkWhere uses regexp to split the link's where clause into smaller parts.
func ParseLinkWhere(str string) string {
	regexpString := `(?ismU)^\s+WHERE\s+`
	parseRegexp := regexp.MustCompile(regexpString)
	result := parseRegexp.ReplaceAll([]byte(str), []byte(""))
	resultAsString := string(result)

	return resultAsString
}

// ParseMapping uses regexp to split the mapping string into smaller parts.
func ParseMapping(str string) (map[string]string, error) {
	result := make(map[string]string)
	regexpString := `(?ismU)^\s*(?P<sourceNode>[^\.,]+)\.(?P<sourceColumn>[^\.,]+)\s+TO\s+(?P<targetNode>[^\.,]+)\.(?P<targetColumn>[^\.,]+)\s*$`
	parseRegexp := regexp.MustCompile(regexpString)
	matches := parseRegexp.FindStringSubmatch(str)
	subNames := parseRegexp.SubexpNames()

	// fmt.Println(matches)

	for i, match := range matches {
		// Skip the first, empty element.
		if i == 0 {
			continue
		}

		// fmt.Println(match)

		result[subNames[i]] = match
	}

	// err := validateMapping(result)

	// return result, err
	return result, nil
}

func validateMapping(result map[string]interface{}) error {
	errorsArr := make([]string, 0)
	var err error = nil

	// fmt.Println(result)

	// entire mapping
	if len(result) == 0 {
		errorsArr = append(errorsArr, "there's a syntax error in the mapping")
	}

	// command
	// if result["command"] == nil || len(result["command"].(string)) == 0 {
	// 	errorsArr = append(errorsArr, "no command found")
	// }

	// // links
	// if result["links"] == nil || len(result["links"].([]map[string]string)) == 0 {
	// 	errorsArr = append(errorsArr, "no links found")
	// }

	if len(errorsArr) > 0 {
		errorsArrJoined := strings.Join(errorsArr, "\n")
		err = &mappingParserError{errMsg: errorsArrJoined}
	}
	return err
}
