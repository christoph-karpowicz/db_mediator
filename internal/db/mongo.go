package db

type MongoDatabase struct {
	db  *DatabaseData
	Tst int
}

func (d *MongoDatabase) Select() string {
	return "mongo"
}
