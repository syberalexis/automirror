package database

import (
	"github.com/boltdb/bolt"
)

// ExistsInDatabase check if a package's name already exists into the database file
func ExistsInDatabase(filename string, name string) (bool, error) {
	db, err := bolt.Open(filename, 0600, nil)
	defer db.Close()
	if err != nil {
		return false, err
	}

	var data []byte
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(filename))
		if err != nil {
			return err
		}
		data = bucket.Get([]byte(name))
		return nil
	})
	return data != nil, err
}

// InsertIntoDatabase insert a package's name into the database file
func InsertIntoDatabase(filename string, key string, value string) error {
	db, err := bolt.Open(filename, 0600, nil)
	defer db.Close()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(filename))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), []byte(value))
	})
}
