package hue

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type authRequest struct {
	DeviceType string `json:"devicetype"`
}

type authResponse struct {
	Error   *errorResp   `json:"error"`
	Success *successResp `json:"success"`
}

type successResp struct {
	Username string `json:"username"`
}

type errorResp struct {
	Type        int    `json:"type"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

func generateUsername() (string, error) {
	// create a buffer for the random seed
	seed := make([]byte, 16)

	// read the seed
	if _, err := rand.Read(seed); err != nil {
		return "", err
	}

	// hash the seed with sha1 for a better distribution
	hashedSeed := sha256.Sum256(seed)

	// return the hex encoded hash
	return hex.EncodeToString(hashedSeed[:10]), nil
}

func authenticate(address string) (string, error) {
	// generate a new username
	id, err := generateUsername()

	if err != nil {
		return "", err
	}

	log.Info("Please press the link button on your bridge")

	// try 30 times to authenticate in intervals if 1 second.
	// This needs to be performed, in order to check, if the
	// link button was pressed
	for i := 0; i < 30; i++ {
		// try to authenticate
		username, err := performAuthRequest(address, id)

		// debug log
		log.Debugf("%v", err)

		// if an error occurred or the link button was not
		// pressed...
		if err != nil {
			// ...wait a second...
			time.Sleep(1 * time.Second)

			// ...and try it again
			continue
		}

		// without an error, the username must be returned
		return username, nil
	}

	// if the authentication failed, an error needs to be
	// returned
	return "", fmt.Errorf("unable to authenticate")
}

func performAuthRequest(address, username string) (string, error) {
	// create a reader for the authentication request body
	bodyBytes, err := json.Marshal(authRequest{
		DeviceType: "HueKit Bridge#" + username,
	})

	// error handling
	if err != nil {
		return "", err
	}

	// create the request
	req, err := http.NewRequest(
		"POST",
		"http://"+address+"/api",
		bytes.NewBuffer(bodyBytes),
	)

	// set the content type to json
	req.Header.Set("Content-Type", "application/json")

	// error handling
	if err != nil {
		return "", err
	}

	// perform the http request
	res, err := http.DefaultClient.Do(req)

	// error handling
	if err != nil {
		return "", err
	}

	// close the body on return in order to avoid
	// memory leaks
	defer res.Body.Close()

	// return the verification result of the response
	return verifyResponse(res.Body)
}

func verifyResponse(res io.ReadCloser) (string, error) {
	// read the complete response body
	resBytes, err := io.ReadAll(res)

	// error handling
	if err != nil {
		return "", err
	}

	// allocate the object, in which the body should be
	// unmarshaled
	var resBody []authResponse

	// unmarshal the body into the previously allocated
	// structure
	err = json.Unmarshal(resBytes, &resBody)

	// error handling
	if err != nil {
		return "", err
	}

	var username string

	// iterate through the results
	for _, result := range resBody {
		if result.Success != nil {
			username = result.Success.Username
		}

		// if the error part is not set, the result is successful...
		if result.Error != nil {
			err = fmt.Errorf(result.Error.Description)
		}
	}

	// indicate a successful authentication
	return username, err
}
