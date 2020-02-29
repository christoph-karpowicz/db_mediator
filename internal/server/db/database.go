package db

import "fmt"

// Database interface is the blueprint for all structs for specific databases.
type Database interface {
	GetConfig() *config
	Init()
	Select(tableName string, conditions string) []map[string]interface{}
	TestConnection()
	Update(table string, key interface{}, column string, val interface{}) (bool, error)
}

// DatabaseError is a custom db error.
type DatabaseError struct {
	DBName string `json:"dbName"`
	ErrMsg string `json:"errMsg"`
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("[database] Database %s %s", e.DBName, e.ErrMsg)
}
