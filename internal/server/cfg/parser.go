package cfg

import (
	"fmt"
	"regexp"
	"strings"

	arrUtil "github.com/christoph-karpowicz/unifier/internal/util/array"
)

type linkParserError struct {
	errMsg string
}

func (e *linkParserError) Error() string {
	return fmt.Sprintf("[link parser] %s", e.errMsg)
}

type mappingParserError struct {
	errMsg string
}

func (e *mappingParserError) Error() string {
	return fmt.Sprintf("[mapping parser] %s", e.errMsg)
}

type matcherParserError struct {
	errMsg string
}

func (e *matcherParserError) Error() string {
	return fmt.Sprintf("['match by' parser] %s", e.errMsg)
}

// ParseLink uses regexp to split the link string into smaller parts.
func ParseLink(link string) (map[string]string, error) {
	result := make(map[string]string)
	ptrn := `(?iU)^\s*` +
		`\[(?P<` + PSUBEXP_SOURCE_NODE + `>[^\.,\s]+)\.(?P<` + PSUBEXP_SOURCE_COLUMN + `>[^\.,\s]+|"[^\.,]+")(\s+)?(?P<` + PSUBEXP_SOURCE_WHERE + `>WHERE\s+[^\s]+.+)?\]` +
		`\s+TO\s+` +
		`\[(?P<` + PSUBEXP_TARGET_NODE + `>[^\.,\s]+)\.(?P<` + PSUBEXP_TARGET_COLUMN + `>[^\.,\s]+|"[^\.,]+")(\s+)?(?P<` + PSUBEXP_TARGET_WHERE + `>WHERE\s+[^\s]+.+)?\]` +
		`\s*$`
	compiledPtrn := regexp.MustCompile(ptrn)
	matches := compiledPtrn.FindStringSubmatch(link)
	subNames := compiledPtrn.SubexpNames()

	if len(matches) == 0 {
		return nil, validateLink(link)
	}

	for i, match := range matches {
		// Skip the first, empty element.
		if i == 0 {
			continue
		}

		if arrUtil.Contains([]string{PSUBEXP_SOURCE_WHERE, PSUBEXP_TARGET_WHERE}, subNames[i]) {
			parsedWhere := ParseLinkWhere(match)
			result[subNames[i]] = parsedWhere
		} else {
			result[subNames[i]] = match
		}
	}
	result["cmd"] = link

	return result, nil
}

func validateLink(link string) error {
	errorsArr := make([]string, 0)
	errorsArr = append(errorsArr, "error in: "+link)
	var err error = nil

	linkTrimmed := strings.Trim(link, " ")

	// Source part of the link.
	sourcePartPtrn := regexp.MustCompile(`^\[.+\].*`)
	sourcePartPtrnMatched := sourcePartPtrn.MatchString(linkTrimmed)
	if !sourcePartPtrnMatched {
		errorsArr = append(errorsArr, "a link has to start with a source in square brackets")
	}
	sourceWherePtrn := regexp.MustCompile(`^\[.+\s+WHERE\s*\]\s+TO.+`)
	if sourcePartPtrnMatched && sourceWherePtrn.MatchString(linkTrimmed) {
		errorsArr = append(errorsArr, "where clause in the source has to be followed by one or more conditions")
	}
	sourceColumnPtrn := regexp.MustCompile(`^\[([^\.,\s]+)\.([^\.,\s]+|"[^\.,]+").*\]\s+TO.+`)
	if sourcePartPtrnMatched && !sourceColumnPtrn.MatchString(linkTrimmed) {
		errorsArr = append(errorsArr, "source node or column name is missing")
	}

	// Middle part of the link.
	middlePartPtrn := regexp.MustCompile(`^\[.+\]\s+TO\s+\[.+\]$`)
	if !middlePartPtrn.MatchString(linkTrimmed) {
		errorsArr = append(errorsArr, "there has to be a 'TO' keyword between the source and target")
	}

	// Target part of the link.
	targetPartPtrn := regexp.MustCompile(`.*\[.+\]$`)
	targetPartPtrnMatched := targetPartPtrn.MatchString(linkTrimmed)
	if !targetPartPtrnMatched {
		errorsArr = append(errorsArr, "a link has to end with a target in square brackets")
	}
	targetWherePtrn := regexp.MustCompile(`.*TO\s+\[.+\s+WHERE\s*\]$`)
	if targetPartPtrnMatched && targetWherePtrn.MatchString(linkTrimmed) {
		errorsArr = append(errorsArr, "where clause in the target has to be followed by one or more conditions")
	}
	targetColumnPtrn := regexp.MustCompile(`.+\[([^\.,\s]+)\.([^\.,\s]+|"[^\.,]+").*\]\s*$`)
	if sourcePartPtrnMatched && !targetColumnPtrn.MatchString(linkTrimmed) {
		errorsArr = append(errorsArr, "target node or column name is missing")
	}

	// no specific erros found
	if len(errorsArr) == 1 {
		errorsArr = append(errorsArr, "there's a syntax error in the link")
	}

	if len(errorsArr) > 1 {
		errorsArrJoined := strings.Join(errorsArr, "\n")
		err = &linkParserError{errMsg: errorsArrJoined}
	}
	return err
}

// ParseLinkWhere uses regexp to split the link's where clause into smaller parts.
func ParseLinkWhere(where string) string {
	ptrn := `(?iU)^\s*WHERE\s+`
	compiledPtrn := regexp.MustCompile(ptrn)
	result := compiledPtrn.ReplaceAll([]byte(where), []byte(""))
	resultAsString := string(result)

	return resultAsString
}

// ParseMapping uses regexp to split the mapping string into smaller parts.
func ParseMapping(mapping string) (map[string]string, error) {
	result := make(map[string]string)
	ptrn := `(?iU)^\s*` +
		`(?P<` + PSUBEXP_SOURCE_NODE + `>[^\.,]+)\.(?P<` + PSUBEXP_SOURCE_COLUMN + `>[^\.,]+)` +
		`\s+TO\s+` +
		`(?P<` + PSUBEXP_TARGET_NODE + `>[^\.,]+)\.(?P<` + PSUBEXP_TARGET_COLUMN + `>[^\.,]+)` +
		`\s*$`
	compiledPtrn := regexp.MustCompile(ptrn)
	matches := compiledPtrn.FindStringSubmatch(mapping)
	subNames := compiledPtrn.SubexpNames()

	if len(matches) == 0 {
		return nil, validateMapping(mapping)
	}

	for i, match := range matches {
		// Skip the first, empty element.
		if i == 0 {
			continue
		}

		result[subNames[i]] = match
	}

	return result, nil
}

func validateMapping(mapping string) error {
	errorsArr := make([]string, 0)
	errorsArr = append(errorsArr, "error in: "+mapping)
	var err error = nil

	mappingTrimmed := strings.Trim(mapping, " ")
	spacePtrn := regexp.MustCompile(`\s+`)
	mappingSplit := spacePtrn.Split(mappingTrimmed, -1)

	if len(mappingSplit) > 2 && mappingSplit[1] != "TO" {
		errorsArr = append(errorsArr, "there has to be a 'TO' keyword between the source and target nodes")
	}
	if len(mappingSplit) == 2 && mappingSplit[0] == "TO" {
		errorsArr = append(errorsArr, "there has to be a source column before the 'TO' keyword")
	}
	if len(mappingSplit) == 2 && mappingSplit[1] == "TO" {
		errorsArr = append(errorsArr, "there has to be a target column after the 'TO' keyword")
	}
	if len(mappingSplit) > 3 {
		errorsArr = append(errorsArr, "there's redundant data in the mapping")
	}
	if len(mappingSplit) < 3 && len(errorsArr) == 1 {
		errorsArr = append(errorsArr, "there's too little data in the mapping")
	}

	// no specific erros found
	if len(errorsArr) == 1 {
		errorsArr = append(errorsArr, "there's a syntax error in the mapping")
	}

	if len(errorsArr) > 1 {
		errorsArrJoined := strings.Join(errorsArr, "\n")
		err = &mappingParserError{errMsg: errorsArrJoined}
	}
	return err
}

// ParseIdsMatcherMethod prepares "ids" method's arguments.
func ParseIdsMatcherMethod(args []string) ([][]string, error) {
	argsSplt := make([][]string, 0)
	for _, arg := range args {
		argSplt := strings.Split(arg, ".")
		argsSplt = append(argsSplt, argSplt)
	}

	validationErr := validateIdsMatcherMethod(args, argsSplt)
	if validationErr != nil {
		return nil, validationErr
	}

	return argsSplt, nil
}

func validateIdsMatcherMethod(args []string, argsSplt [][]string) error {
	errorsArr := make([]string, 0)
	var err error = nil

	if len(args) > 2 {
		errorsArr = append(errorsArr, "too many arguments given for this match method")
	}
	if len(args) < 2 {
		errorsArr = append(errorsArr, "too few arguments given for this match method")
	}
	if len(args) == 2 && (len(argsSplt[0]) != 2 || len(argsSplt[1]) != 2) {
		errorsArr = append(errorsArr, "each argument has to consist of node name and ID column name separated by a dot")
	}
	if len(errorsArr) == 0 && argsSplt[0][0] == argsSplt[1][0] {
		errorsArr = append(errorsArr, "\"ids\" match method accepts only external ID column names from different nodes")
	}

	if len(errorsArr) > 0 {
		errorsArrJoined := strings.Join(errorsArr, "\n")
		err = &matcherParserError{errMsg: errorsArrJoined}
	}
	return err
}
