// Package hue Implements access to the hue bridge
package hue

import (
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/dj95/huekit/pkg/store"
)

var modelIDPattern *regexp.Regexp = nil

// Bridger Interface for interacting with the hue bridge
type Bridger interface {
	Lights() ([]*Light, error)
}

// Bridge Implements handling with the hue bridge
type Bridge struct {
	address  string
	username string
}

// NewBridge Instantiates a new bridge with the given store. If no
// authentication is saved, it will authenticate against the bridge
func NewBridge(address string, store store.Store) (Bridger, error) {
	// check if the username is already set in the database
	username, err := store.Get("bridge_username")

	log.Debugf("%v ˆ%v", username, err)

	// handle the error, if the username does not exist
	if err != nil {
		// authenticate
		username, err = authenticate(address)
	}

	// handle authentication error
	if err != nil {
		return nil, err
	}

	// update the username in the database
	err = store.Set("bridge_username", username)

	// error handling
	if err != nil {
		return nil, err
	}

	// return the initialized bridge
	return &Bridge{
		address:  address,
		username: username,
	}, nil
}

// ModelIDIsFromHue Check if the modelID matches the pattern of
// a hue model id or not.
func ModelIDIsFromHue(modelID string) bool {
	// check if the pattern in this package is initialized
	if modelIDPattern == nil {
		// if not, initialize it with a compiled pattern in order to
		// gain a better performance  in contrast to compiling it
		// every time, this method is called
		modelIDPattern = regexp.MustCompile(`^[A-Z]{3}[0-9]{3}$`)
	}

	// return if the pattern matches or not
	return modelIDPattern.MatchString(modelID)
}
