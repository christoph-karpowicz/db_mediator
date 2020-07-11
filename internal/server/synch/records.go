package synch

import "reflect"

type records []*record

func (tr records) FindRecordPointer(searchedRecord map[string]interface{}) *record {
	for i := range tr {
		if reflect.DeepEqual(tr[i].Data, searchedRecord) {
			return tr[i]
		}
	}
	return nil
}
