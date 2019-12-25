package synch

import "reflect"

type tableRecords struct {
	table   *table
	records []record
}

func (tr tableRecords) FindRecordPointer(searchedRecord map[string]interface{}) *record {
	for _, rec := range tr.records {
		if reflect.DeepEqual(rec.Data, searchedRecord) {
			return &rec
		}
	}
	return nil
}
