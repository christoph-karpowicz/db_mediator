package db

type UpdateDto struct {
	TableName         string
	KeyName           string
	KeyValue          interface{}
	UpdatedColumnName string
	NewValue          interface{}
}

type InsertDto struct {
	TableName string
	KeyName   string
	KeyValue  interface{}
	Values    map[string]interface{}
}
