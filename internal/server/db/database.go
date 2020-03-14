package db

import "fmt"

// Database interface is the blueprint for all structs for specific databases.
type Database interface {
	GetConfig() *config
	Init()
	Select(tableName string, conditions string) []map[string]interface{}
	TestConnection()
	Update(table string, keyName string, keyVal interface{}, column string, val interface{}) (bool, error)
	Insert(table string, keyName string, keyVal interface{}, values map[string]interface{}) (bool, error)
}

// DatabaseError is a custom db error.
type DatabaseError struct {
	DBName  string      `json:"dbName"`
	ErrMsg  string      `json:"errMsg"`
	KeyName string      `json:"keyName"`
	KeyVal  interface{} `json:"keyVal"`
}

func (e *DatabaseError) Error() string {
	if e.KeyName != "" && e.KeyVal != nil {
		return fmt.Sprintf("[database %s] %s (key: %s, val: %v).", e.DBName, e.ErrMsg, e.KeyName, e.KeyVal)
	}

	return fmt.Sprintf("[database %s] %s", e.DBName, e.ErrMsg)
}
