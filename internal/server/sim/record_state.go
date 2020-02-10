package sim

// RecordState holds information about a certain state of a data record.
type RecordState struct {
	KeyName      interface{}
	KeyValue     interface{}
	ColumnName   string
	CurrentValue interface{}
	NewValue     interface{}
}
