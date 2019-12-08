package db

import (
	"database/sql"
	"fmt"
)

type PostgresDatabase struct {
	DB               *DatabaseData
	connectionString string
}

func (d *PostgresDatabase) GetData() *DatabaseData {
	return d.DB
}

func (d *PostgresDatabase) Init() {
	d.connectionString = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		d.DB.Host, d.DB.Port, d.DB.User, d.DB.Password, d.DB.Name)
}

func (d *PostgresDatabase) SelectAll(tableName string) []map[string]interface{} {
	d.TestConnection()

	database, err := sql.Open("postgres", d.connectionString)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	rows, err := database.Query(fmt.Sprintf(`SELECT * FROM %s WHERE title ILIKE 'Des%%'`, tableName))
	if err != nil {
		panic(err)
	}

	cols, _ := rows.Columns()

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			panic(err)
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		fmt.Println(m)
	}
	return nil
}

func (d *PostgresDatabase) TestConnection() {
	database, err := sql.Open("postgres", d.connectionString)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}
