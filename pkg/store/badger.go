package store

import (
	badger "github.com/dgraph-io/badger/v2"
)

// Badger Implement the key-value storage with badger
type Badger struct {
	db *badger.DB
}

// NewBadger Instatiate a new store with a badger db
func NewBadger(db *badger.DB) Store {
	return &Badger{
		db: db,
	}
}

// Get Retrieve a key from the badger db
func (b *Badger) Get(key string) (string, error) {
	var result string

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))

		if err != nil {
			return err
		}

		item.Value(func(value []byte) error {
			result = string(value)

			return nil
		})

		return err
	})

	return result, err
}

// Set Save a key value pair into the badger db
func (b *Badger) Set(key, value string) error {
	// run a db update callback function
	return b.db.Update(func(txn *badger.Txn) error {
		// set the key to the value
		return txn.Set([]byte(key), []byte(value))
	})
}
