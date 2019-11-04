package utils

import "database/sql"

func InitializeDatabase(filename string, query string) error {
	database, err := sql.Open("sqlite3", filename+"?cache=shared")
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
	database, err := sql.Open("sqlite3", filename+"?cache=shared")
	if err != nil {
		return false, err
	}

	statement, err := database.Prepare(query)
	if err != nil {
		return false, err
	}

	rows, err := statement.Query(args...)
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func InsertIntoDatabase(filename string, query string, args ...interface{}) error {
	database, err := sql.Open("sqlite3", filename+"?cache=shared")
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
