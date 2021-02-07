package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/christoph-karpowicz/db_mediator/internal/server/cfg"
	_ "github.com/lib/pq"
)

// PostgresDatabase implements Database interface for PostgreSQL database.
type postgresDatabase struct {
	cfg              *cfg.DbConfig
	connectionString string
}

// GetConfig returns information about the database, which was parsed from JSON.
func (d *postgresDatabase) GetConfig() *cfg.DbConfig {
	return d.cfg
}

// Init creates the db connection string.
func (d *postgresDatabase) Init() {
	d.connectionString = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		d.cfg.Host, d.cfg.Port, d.cfg.User, d.cfg.Password, d.cfg.Name)

	d.TestConnection()
}

// Insert inserts one row into a given table.
func (d *postgresDatabase) Insert(inDto InsertDto) error {
	database, err := sql.Open("postgres", d.connectionString)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	var columnList []string = make([]string, 0)
	var valuesList []interface{} = make([]interface{}, 0)
	var valuesPlaceholderList []string = make([]string, 0)
	var valuesCounter int64 = 1

	for key, val := range inDto.Values {
		valuesCounterStr := strconv.FormatInt(valuesCounter, 10)

		columnList = append(columnList, key)
		valuesList = append(valuesList, val)
		valuesPlaceholderList = append(valuesPlaceholderList, "$"+valuesCounterStr)
		valuesCounter++
	}

	query := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", inDto.TableName, strings.Join(columnList, ", "), strings.Join(valuesPlaceholderList, ", "))

	result, err := database.Exec(query, valuesList...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		dbErr := &DatabaseError{DBName: d.cfg.Name, ErrMsg: "row hasn't been inserted" /* , KeyName: keyName, KeyValue: keyVal */}
		return dbErr
	}

	return nil
}

// Select selects data from the database, with or without a WHERE clause.
func (d *postgresDatabase) Select(tableName string, conditions string) []map[string]interface{} {
	var allRecords []map[string]interface{}

	database, err := sql.Open("postgres", d.connectionString)
	if err != nil {
		panic(&DatabaseError{DBName: d.cfg.Name, ErrMsg: err.Error()})
	}
	defer database.Close()

	if conditions != "" {
		conditions = fmt.Sprintf(" WHERE %s", conditions)
	}

	query := fmt.Sprintf("SELECT * FROM %s%s", tableName, conditions)

	rows, err := database.Query(query)
	if err != nil {
		panic(&DatabaseError{DBName: d.cfg.Name, ErrMsg: err.Error()})
	}

	cols, _ := rows.Columns()

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers.
		if err := rows.Scan(columnPointers...); err != nil {
			panic(&DatabaseError{DBName: d.cfg.Name, ErrMsg: err.Error()})
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		record := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			record[colName] = *val
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		allRecords = append(allRecords, record)
	}

	return allRecords
}

// TestConnection pings the database.
func (d *postgresDatabase) TestConnection() {
	database, err := sql.Open("postgres", d.connectionString)
	if err != nil {
		panic(&DatabaseError{DBName: d.cfg.Name, ErrMsg: err.Error()})
	}
	defer database.Close()

	err = database.Ping()
	if err != nil {
		panic(&DatabaseError{DBName: d.cfg.Name, ErrMsg: err.Error()})
	}

	fmt.Println("Successfully connected!")
}

// Update updates a record with the provided key.
func (d *postgresDatabase) Update(upDto UpdateDto) error {
	database, err := sql.Open("postgres", d.connectionString)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	query := fmt.Sprintf("UPDATE %s SET %s = $1 WHERE %s = $2", upDto.TableName, upDto.UpdatedColumnName, upDto.KeyName)

	result, err := database.Exec(query, upDto.NewValue, upDto.KeyValue)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		dbErr := &DatabaseError{DBName: d.cfg.Name, ErrMsg: "no rows affected in update", KeyName: upDto.KeyName, KeyValue: upDto.KeyValue}
		return dbErr
	}

	return nil
}
