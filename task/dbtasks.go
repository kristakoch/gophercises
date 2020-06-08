package main

import (
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

// OpenDB opens the database.
func OpenDB() error { // take in the path to the db (tasks.db)
	var err error
	db, err = bolt.Open("tasks.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("TaskBucket"))
		return err
	})
}

// AddToDB adds an entry to the database.
func AddToDB(task string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TaskBucket"))
		err := b.Put([]byte(task), []byte(""))
		return err
	})
	return err
}
