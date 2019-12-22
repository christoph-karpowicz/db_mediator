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

func (d *PostgresDatabase) Select(tableName string, conditions string) []map[string]interface{} {
	var allRecords []map[string]interface{}

	d.TestConnection()

	database, err := sql.Open("postgres", d.connectionString)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	if conditions != "-" && conditions != "" {
		conditions = fmt.Sprintf(` WHERE %s`, conditions)
	} else {
		conditions = ""
	}

	query := fmt.Sprintf(`SELECT * FROM %s%s`, tableName, conditions)

	rows, err := database.Query(query)
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
		record := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			record[colName] = *val
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		// fmt.Println(record)
		allRecords = append(allRecords, record)
	}

	return allRecords
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
