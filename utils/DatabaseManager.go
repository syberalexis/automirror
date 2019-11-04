package utils

import (
	"database/sql"
	"fmt"
)

func InitializeDatabase(filename string, query string) error {
	database, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared", filename))
	if err != nil {
		return err
	}
	statement, err := database.Prepare(query)
	if err != nil {
		return err
	}

	_, err = statement.Exec()
	return err
}

func ExistsInDatabase(filename string, query string, args ...interface{}) (bool, error) {
	database, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared", filename))
	if err != nil {
		return false, err
	}

	statement, err := database.Prepare(query)
	if err != nil {
		return false, err
	}

	rows, err := statement.Query(args...)
	defer rows.Close()
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func InsertIntoDatabase(filename string, query string, args ...interface{}) error {
	database, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared", filename))
	if err != nil {
		return err
	}
	statement, err := database.Prepare(query)
	if err != nil {
		return err
	}

	_, err = statement.Exec(args...)
	return err
}
