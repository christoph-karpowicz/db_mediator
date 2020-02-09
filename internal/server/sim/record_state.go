package sim

type RecordState struct {
	keyName      interface{}
	keyValue     interface{}
	columnName   string
	currentValue interface{}
	newValue     interface{}
}
