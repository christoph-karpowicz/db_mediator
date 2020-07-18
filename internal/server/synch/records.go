package synch

import (
	"errors"
	"reflect"
)

type records []*record

func (tr records) FindRecordPointer(searchedRecord map[string]interface{}) (*record, error) {
	for i := range tr {
		if reflect.DeepEqual(tr[i].Data, searchedRecord) {
			return tr[i], nil
		}
	}
	return nil, errors.New("[data selection] Record hasn't been found")
}
