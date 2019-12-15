package synch

import "reflect"

type TableRecords struct {
	records []Record
}

func (tr TableRecords) FindRecordPointer(searchedRecord map[string]interface{}) *Record {
	for _, rec := range tr.records {
		if reflect.DeepEqual(rec.Data, searchedRecord) {
			return &rec
		}
	}
	return nil
}
