package db

type Database interface {
	GetData() *DatabaseData
	Select() string
}
