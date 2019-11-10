package db

type PostgresDatabase struct {
	db *DatabaseData
}

func (d *PostgresDatabase) Select() string {
	return "postgres"
}
