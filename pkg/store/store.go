// Package store Store key value pairs
package store

// Store Interface for saving and retrieving key-value pairs
type Store interface {
	// Get Return a value related to the key from the store
	Get(key string) (string, error)

	// Set Saves a key-value relation to the store
	Set(key, value string) error
}
