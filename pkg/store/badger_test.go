package store

import (
	"testing"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/stretchr/testify/assert"
)

func prepareDB(data map[string]string) (*badger.DB, error) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))

	if err != nil {
		return nil, err
	}

	for key, value := range data {
		err := db.Update(func(txn *badger.Txn) error {
			// set the key to the value
			return txn.Set([]byte(key), []byte(value))
		})

		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func TestBadger_Get(t *testing.T) {
	tests := []struct {
		description    string
		dbData         map[string]string
		key            string
		expectedError  bool
		expectedResult string
	}{
		{
			description: "existing key",
			dbData: map[string]string{
				"foo": "bar",
			},
			key:            "foo",
			expectedError:  false,
			expectedResult: "bar",
		},
	}

	for _, test := range tests {
		// create an in memory db with required data
		db, err := prepareDB(test.dbData)

		// on db creation error -> fail
		assert.Nilf(t, err, test.description)

		// create the badger store and get the key
		result, err := NewBadger(db).Get(test.key)

		// assert the expected behaviour
		assert.Equalf(t, test.expectedError, err != nil, test.description)
		assert.Equalf(t, test.expectedResult, result, test.description)
	}
}
