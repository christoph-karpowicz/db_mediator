package db

type Database interface {
	GetData() *DatabaseData
	Init()
	SelectAll(tableName string) []map[string]interface{}
	TestConnection()
}
