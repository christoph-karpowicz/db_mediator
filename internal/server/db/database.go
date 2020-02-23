package db

// Database interface is the blueprint for all structs for specific databases.
type Database interface {
	GetConfig() *config
	Init()
	Select(tableName string, conditions string) []map[string]interface{}
	TestConnection()
	Update(table string, key interface{}, column string, val interface{}) (bool, error)
}
