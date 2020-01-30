package validationUtil

import (
	"reflect"
	"strconv"
	"strings"

	arrayUtil "github.com/christoph-karpowicz/unifier/internal/util/array"
)

func YAMLField(fieldValue interface{}, fieldName string) bool {
	switch fieldValue.(type) {
	case int:
		if fieldValue.(int) <= 0 {
			return false
		}
	case string:
		if len(fieldValue.(string)) <= 0 {
			return false
		}
	case []string:
		if len(fieldValue.([]string)) == 0 {
			return false
		}
		for i, val := range fieldValue.([]string) {
			if !YAMLField(val, strconv.Itoa(i)) {
				return false
			}
		}
	}
	return true
}

func YAMLStruct(structure interface{}, nullableFields []string) {
	fieldValue := reflect.ValueOf(structure)
	fieldType := fieldValue.Type()

	for i := 0; i < fieldValue.NumField(); i++ {

		// If a field doesn't have a tag, it means it's value wasn't imported from YAML file.
		if fieldType.Field(i).Tag == "" {
			continue
		}

		if !YAMLField(fieldValue.Field(i).Interface(), fieldType.Field(i).Name) {
			if !arrayUtil.Contains(nullableFields, strings.ToLower(fieldType.Field(i).Name)) {
				panic(fieldType.Field(i).Name + " is invalid.")
			}
		}

	}
}
