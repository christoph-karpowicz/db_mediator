package db

type MongoDatabase struct {
	DB  *DatabaseData
	Tst int
}

func (d *MongoDatabase) GetData() *DatabaseData {
	return d.DB
}

func (d *MongoDatabase) Select() string {
	return "mongo"
}
