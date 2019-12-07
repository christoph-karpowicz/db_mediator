package validationUtil

import (
	"reflect"
	"strconv"
)

func JSONField(fieldValue interface{}, fieldName string) string {
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
			if JSONField(val, strconv.Itoa(i)) != "" {
				return fieldName
			}
		}
	}
	return invalid
}

func JSONStruct(structure interface{}) {
	fieldValue := reflect.ValueOf(structure)
	fieldType := fieldValue.Type()

	for i := 0; i < fieldValue.NumField(); i++ {

		if invalidField := JSONField(fieldValue.Field(i).Interface(), fieldType.Field(i).Name); invalidField != "" {
			panic(invalidField + " is invalid.")
		}

	}
}
