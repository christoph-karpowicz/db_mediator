package db

type PostgresDatabase struct {
	DB *DatabaseData
}

func (d *PostgresDatabase) GetData() *DatabaseData {
	return d.DB
}

func (d *PostgresDatabase) Select() string {
	return "postgres"
}
