package lang

import (
	"regexp"
)

// ParseMapping uses regexp to split the mapping string into smaller parts.
func ParseMapping(str string) map[string]interface{} {
	rawMpng := make(map[string]interface{})
	r := regexp.MustCompile(`(?isU)^\s*(?P<command>[a-z]+)\s+(?P<links>\[.*\]\sTO\s\[.*\],?)+\s+MATCH\sBY\s(?P<matchBy>[a-z]+\(.+\))\s+DO\s(?P<do>[a-z]+)\s*$`)
	matches := r.FindStringSubmatch(str)
	subNames := r.SubexpNames()

	for i, match := range matches {
		// Skip the first, empty element.
		if i == 0 {
			continue
		}

		if subNames[i] == "links" {
			rawMpng[subNames[i]] = make([]map[string]string, 0)
			for _, link := range regexp.MustCompile(`(?s)\s*,\s*`).Split(match, -1) {
				rawMpng[subNames[i]] = append(rawMpng[subNames[i]].([]map[string]string), parseLink(link))
				// fmt.Println(rawMpng[subNames[i]])
			}
		} else {
			rawMpng[subNames[i]] = match
		}
	}

	return rawMpng
}

// parseLink splits an individual link into smaller parts.
func parseLink(str string) map[string]string {
	parsedLink := make(map[string]string)
	r := regexp.MustCompile(`(?iU)^\[(?P<sourceNode>.+)\.(?P<sourceColumn>.+)(\s+WHERE\s+(?P<sourceWhere>.+))?\]\sTO\s\[(?P<targetNode>.+)\.(?P<targetColumn>.+)(\s+WHERE\s+(?P<targetWhere>.+))?\]$`)
	matches := r.FindStringSubmatch(str)
	subNames := r.SubexpNames()

	for i, match := range matches {
		if subNames[i] != "" {
			parsedLink[subNames[i]] = match
		}
	}

	return parsedLink
}
