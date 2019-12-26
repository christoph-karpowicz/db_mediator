package validationUtil

import (
	"reflect"
	"strconv"
)

func YAMLField(fieldValue interface{}, fieldName string) string {
	var invalid string
	switch fieldValue.(type) {
	case int:
		if fieldValue.(int) <= 0 {
			invalid = fieldName
		}
	case string:
		if len(fieldValue.(string)) <= 0 {
			invalid = fieldName
		}
	case []string:
		if len(fieldValue.([]string)) < 2 {
			return fieldName
		}
		for i, val := range fieldValue.([]string) {
			if YAMLField(val, strconv.Itoa(i)) != "" {
				return fieldName
			}
		}
	}
	return invalid
}

func YAMLStruct(structure interface{}) {
	fieldValue := reflect.ValueOf(structure)
	fieldType := fieldValue.Type()

	for i := 0; i < fieldValue.NumField(); i++ {

		if invalidField := YAMLField(fieldValue.Field(i).Interface(), fieldType.Field(i).Name); invalidField != "" {
			panic(invalidField + " is invalid.")
		}

	}
}
