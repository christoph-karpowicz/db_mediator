package synch

import (
	"errors"
	"reflect"
)

type records []*record

func (r records) findRecordPointer(searchedRecord map[string]interface{}) (*record, error) {
	for i := range r {
		if reflect.DeepEqual(r[i].Data, searchedRecord) {
			return r[i], nil
		}
	}
	return nil, errors.New("[data selection] Record hasn't been found")
}

func (r *records) setActiveIn(lnk *Link) {
	for _, record := range *r {
		record.ActiveIn = append(record.ActiveIn, lnk)
	}
}
