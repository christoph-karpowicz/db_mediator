package synch

import (
	"errors"
	"reflect"
)

func AreEqual(val1 interface{}, val2 interface{}) (bool, error) {
	var val1kind reflect.Kind = reflect.TypeOf(val1).Kind()
	var val2kind reflect.Kind = reflect.TypeOf(val2).Kind()

	// Strings
	if val1kind == reflect.String && val2kind == reflect.String {
		return val1.(string) == val2.(string), nil
	}

	// Signed ints
	if isSignedInt(val1kind) && isSignedInt(val2kind) {
		var val1int64 int64
		var val2int64 int64

		switch val1kind {
		case reflect.Int:
			val1int64 = int64(val1.(int))
		case reflect.Int8:
			val1int64 = int64(val1.(int8))
		case reflect.Int16:
			val1int64 = int64(val1.(int16))
		case reflect.Int32:
			val1int64 = int64(val1.(int32))
		case reflect.Int64:
			val1int64 = int64(val1.(int64))
		}

		switch val2kind {
		case reflect.Int:
			val2int64 = int64(val2.(int))
		case reflect.Int8:
			val2int64 = int64(val2.(int8))
		case reflect.Int16:
			val2int64 = int64(val2.(int16))
		case reflect.Int32:
			val2int64 = int64(val2.(int32))
		case reflect.Int64:
			val2int64 = int64(val2.(int64))
		}

		return val1int64 == val2int64, nil
	}

	// Unsigned ints
	if isUnsignedInt(val1kind) && isUnsignedInt(val2kind) {
		var val1uint64 uint64
		var val2uint64 uint64

		switch val1kind {
		case reflect.Uint:
			val1uint64 = uint64(val1.(uint))
		case reflect.Uint8:
			val1uint64 = uint64(val1.(uint8))
		case reflect.Uint16:
			val1uint64 = uint64(val1.(uint16))
		case reflect.Uint32:
			val1uint64 = uint64(val1.(uint32))
		case reflect.Uint64:
			val1uint64 = uint64(val1.(uint64))
		}

		switch val2kind {
		case reflect.Uint:
			val2uint64 = uint64(val2.(uint))
		case reflect.Uint8:
			val2uint64 = uint64(val2.(uint8))
		case reflect.Uint16:
			val2uint64 = uint64(val2.(uint16))
		case reflect.Uint32:
			val2uint64 = uint64(val2.(uint32))
		case reflect.Uint64:
			val2uint64 = uint64(val2.(uint64))
		}

		return val1uint64 == val2uint64, nil
	}

	return false, errors.New("Invalid data types.")
}

func isSignedInt(val reflect.Kind) bool {
	signedIntTypes := []reflect.Kind{
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
	}
	for _, t := range signedIntTypes {
		if val == t {
			return true
		}
	}
	return false
}

func isUnsignedInt(val reflect.Kind) bool {
	unsignedIntTypes := []reflect.Kind{
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
	}
	for _, t := range unsignedIntTypes {
		if val == t {
			return true
		}
	}
	return false
}

func MapToRecords(mapRecordArray []map[string]interface{}) []Record {
	var recordArray []Record = make([]Record, 0)
	for _, mapRecord := range mapRecordArray {
		record := Record{Data: mapRecord, IsActive: false}
		recordArray = append(recordArray, record)
	}
	return recordArray
}
