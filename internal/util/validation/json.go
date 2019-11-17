package validationUtil

import (
	"reflect"
	"strconv"
)

func JSONField(fieldValue interface{}, fieldName string, ignore string) string {
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
		if len(fieldValue.([]string)) < 2 && fieldName != ignore {
			return fieldName
		}
		for i, val := range fieldValue.([]string) {
			if JSONField(val, strconv.Itoa(i), "") != "" {
				return fieldName
			}
		}
	}
	return invalid
}

func JSONStruct(structure interface{}, ignore string) {
	fieldValue := reflect.ValueOf(structure)
	fieldType := fieldValue.Type()

	for i := 0; i < fieldValue.NumField(); i++ {

		if invalidField := JSONField(fieldValue.Field(i).Interface(), fieldType.Field(i).Name, ignore); invalidField != "" {
			panic(invalidField + " is invalid.")
		}

	}
}
