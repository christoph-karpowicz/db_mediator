package synch

import "reflect"

type tableRecords struct {
	table   *table
	records []record
}

func (tr tableRecords) FindRecordPointer(searchedRecord map[string]interface{}) *record {
	for i := range tr.records {
		if reflect.DeepEqual(tr.records[i].Data, searchedRecord) {
			return &tr.records[i]
		}
	}
	return nil
}
