package utils

import (
	"database/sql"
	"fmt"
)

func InitializeDatabase(filename string) error {
	database, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared", filename))
	defer database.Close()
	if err != nil {
		return err
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS package (id INTEGER PRIMARY KEY, `name` TEXT)")
	defer statement.Close()
	if err != nil {
		return err
	}

	_, err = statement.Exec()
	return err
}

func ExistsInDatabase(filename string, name string) (bool, error) {
	database, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared", filename))
	defer database.Close()
	if err != nil {
		return false, err
	}

	statement, err := database.Prepare("SELECT id FROM package WHERE `name` = ?")
	defer statement.Close()
	if err != nil {
		return false, err
	}

	rows, err := statement.Query(name)
	defer rows.Close()
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func InsertIntoDatabase(filename string, name string) error {
	database, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared", filename))
	defer database.Close()
	if err != nil {
		return err
	}
	statement, err := database.Prepare("INSERT INTO package (`name`) VALUES (?)")
	defer statement.Close()
	if err != nil {
		return err
	}

	_, err = statement.Exec(name)
	return err
}
