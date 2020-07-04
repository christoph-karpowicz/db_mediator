package cfg

import (
	"fmt"
	"regexp"
	"strings"

	arrUtil "github.com/christoph-karpowicz/unifier/internal/util/array"
)

type mappingParserError struct {
	errMsg string
}

func (e *mappingParserError) Error() string {
	return fmt.Sprintf("[mapping parser] %s", e.errMsg)
}

// ParseLink uses regexp to split the link string into smaller parts.
func ParseLink(link string) (map[string]string, error) {
	result := make(map[string]string)
	ptrn := `(?isU)^\s*\` +
		`[(?P<sourceNode>[^\.,]+)\.(?P<sourceColumn>[^\.,]+)(?P<sourceWhere>\s+WHERE\s+.+)?\]` +
		`\s+TO\s+` +
		`\[(?P<targetNode>[^\.,]+)\.(?P<targetColumn>[^\.,]+)(?P<targetWhere>\s+WHERE\s+.+)?\]` +
		`\s*$`
	compiledPtrn := regexp.MustCompile(ptrn)
	matches := compiledPtrn.FindStringSubmatch(link)
	subNames := compiledPtrn.SubexpNames()

	// fmt.Println(matches)

	for i, match := range matches {
		// Skip the first, empty element.
		if i == 0 {
			continue
		}

		// fmt.Println(match)

		if arrUtil.Contains([]string{"sourceWhere", "targetWhere"}, subNames[i]) {
			parsedWhere := ParseLinkWhere(match)
			result[subNames[i]] = parsedWhere
		} else {
			result[subNames[i]] = match
		}
	}

	result["cmd"] = link

	// err := validateMapping(result)

	// return result, err
	return result, nil
}

// ParseLinkWhere uses regexp to split the link's where clause into smaller parts.
func ParseLinkWhere(where string) string {
	ptrn := `(?isU)^\s+WHERE\s+`
	compiledPtrn := regexp.MustCompile(ptrn)
	result := compiledPtrn.ReplaceAll([]byte(where), []byte(""))
	resultAsString := string(result)

	return resultAsString
}

// ParseMapping uses regexp to split the mapping string into smaller parts.
func ParseMapping(mapping string) (map[string]string, error) {
	result := make(map[string]string)
	ptrn := `(?isU)^\s*` +
		`(?P<sourceNode>[^\.,]+)\.(?P<sourceColumn>[^\.,]+)` +
		`\s+TO\s+` +
		`(?P<targetNode>[^\.,]+)\.(?P<targetColumn>[^\.,]+)` +
		`\s*$`
	compiledPtrn := regexp.MustCompile(ptrn)
	matches := compiledPtrn.FindStringSubmatch(mapping)
	subNames := compiledPtrn.SubexpNames()

	if len(matches) == 0 {
		return nil, validateMapping(mapping)
	}
	// fmt.Println(matches)

	for i, match := range matches {
		// Skip the first, empty element.
		if i == 0 {
			continue
		}

		// fmt.Println(match)

		result[subNames[i]] = match
	}

	return result, nil
}

func validateMapping(mapping string) error {
	errorsArr := make([]string, 0)
	var err error = nil

	spacePtrn := regexp.MustCompile(`\s+`)
	mappingSplit := spacePtrn.Split(mapping, 3)
	fmt.Println(mappingSplit)
	// if !compiledToKeywordPtrn.MatchString(mapping) {
	// 	errorsArr = append(errorsArr, "there has to be a 'TO' keyword between the source and target nodes")
	// }

	// no specific erros found
	if len(errorsArr) == 0 {
		errorsArr = append(errorsArr, "there's a syntax error in the mapping")
	}

	if len(errorsArr) > 0 {
		errorsArrJoined := strings.Join(errorsArr, "\n")
		err = &mappingParserError{errMsg: errorsArrJoined}
	}
	return err
}
