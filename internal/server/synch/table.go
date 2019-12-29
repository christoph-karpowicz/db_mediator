package synch

type table struct {
	id         string
	dbName     string
	name       string
	oldRecords *tableRecords
	records    *tableRecords
}
