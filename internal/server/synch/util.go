package synch

func MapToRecords(mapRecordArray []map[string]interface{}) []Record {
	var recordArray []Record = make([]Record, 0)
	for _, mapRecord := range mapRecordArray {
		record := Record{Data: mapRecord, IsActive: false}
		recordArray = append(recordArray, record)
	}
	return recordArray
}
