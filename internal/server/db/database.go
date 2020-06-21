package db

import "fmt"

// Database interface is the blueprint for all structs for specific databases.
type Database interface {
	GetConfig() *config
	Init()
	Select(tableName string, conditions string) []map[string]interface{}
	TestConnection()
	Insert(inDto InsertDto) (bool, error)
	Update(upDto UpdateDto) (bool, error)
}

// DatabaseError is a custom db error.
type DatabaseError struct {
	DBName   string      `json:"dbName"`
	ErrMsg   string      `json:"errMsg"`
	KeyName  string      `json:"keyName"`
	KeyValue interface{} `json:"keyVal"`
}

func (e *DatabaseError) Error() string {
	if e.KeyName != "" && e.KeyValue != nil {
		return fmt.Sprintf("[database %s] %s (key: %s, val: %v).", e.DBName, e.ErrMsg, e.KeyName, e.KeyValue)
	}

	return fmt.Sprintf("[database %s] %s", e.DBName, e.ErrMsg)
}
