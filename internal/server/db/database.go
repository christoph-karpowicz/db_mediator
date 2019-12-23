package db

type Database interface {
	GetData() *DatabaseData
	Init()
	Select(tableName string, conditions string) []map[string]interface{}
	TestConnection()
	Update(key interface{}, val interface{}) (bool, error)
}
