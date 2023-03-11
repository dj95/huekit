package hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Light Represents a light/plug at the hue bridge
type Light struct {
	ID               string
	Type             string            `json:"type"`
	Name             string            `json:"name"`
	ModelID          string            `json:"modelid"`
	ManufacturerName string            `json:"manufacturername"`
	SoftwareVersion  string            `json:"swversion"`
	State            *State            `json:"state"`
	PointSymbol      map[string]string `json:"pointsymbol"`
}

// State Represents the state of a light
type State struct {
	On               bool      `json:"on"`
	Brightness       int       `json:"bri,omitempty"`
	Hue              int       `json:"hue,omitempty"`
	Saturation       int       `json:"sat,omitempty"`
	XY               []float64 `json:"xy,omitempty"`
	ColorTemperature int       `json:"ct,omitempty"`
	Alert            string    `json:"alert,omitempty"`
	Effect           string    `json:"effect,omitempty"`
	ColorMode        string    `json:"colormode,omitempty"`
	Reachable        bool      `json:"reachable,omitempty"`
}

// LightName Represents the lights name in the /lights api call
type LightName struct {
	Name string `json:"name"`
}

// Lights Query and return all lights
func (b *Bridge) Lights() ([]*Light, error) {
	// perform the api request to fetch all lights
	res, err := http.Get(
		"http://" + b.address + "/api/" + b.username + "/lights",
	)

	// handle http errors
	if err != nil {
		return nil, err
	}

	// close the response body on return in order to avoid memory
	// leaks
	defer res.Body.Close()

	// allocate the structure for the response body in memory
	var lightIDs map[string]*LightName

	// read the body
	bodyBytes, err := io.ReadAll(res.Body)

	// handle read errors
	if err != nil {
		return nil, err
	}

	// unmarshal the json body
	err = json.Unmarshal(bodyBytes, &lightIDs)

	// handle json decoding errors
	if err != nil {
		return nil, err
	}

	var lights []*Light

	for id := range lightIDs {
		light, err := b.Light(id)

		if err != nil {
			return nil, err
		}

		lights = append(lights, light)
	}

	// return the result
	return lights, nil
}

// Light Query and return a light by its id
func (b *Bridge) Light(id string) (*Light, error) {
	// perform the api request to fetch all lights
	res, err := http.Get(
		"http://" + b.address + "/api/" + b.username + "/lights/" + id,
	)

	// handle http errors
	if err != nil {
		return nil, err
	}

	// close the response body on return in order to avoid memory
	// leaks
	defer res.Body.Close()

	// allocate the structure for the response body in memory
	var light Light

	// read the body
	bodyBytes, err := io.ReadAll(res.Body)

	// handle read errors
	if err != nil {
		return nil, err
	}

	// unmarshal the json body
	err = json.Unmarshal(bodyBytes, &light)

	// handle json decoding errors
	if err != nil {
		return nil, err
	}

	// add the ID to the light
	light.ID = id

	return &light, nil
}

type toggleResponse struct {
	Error   *errorResp   `json:"error"`
	Success *successResp `json:"success"`
}

// LightUpdateState Update the state of a light
func (b *Bridge) LightUpdateState(light *Light, state *State) error {
	// create the request body
	body, err := json.Marshal(state)

	if err != nil {
		return err
	}

	// perform the api request to fetch all lights
	req, err := http.NewRequest(
		"PUT",
		"http://"+b.address+"/api/"+b.username+"/lights/"+light.ID+"/state",
		bytes.NewBuffer(body),
	)

	// handle http errors
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	resByte, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	var toggleResp []toggleResponse

	if err := json.Unmarshal(resByte, &toggleResp); err != nil {
		return err
	}

	err = nil

	for _, res := range toggleResp {
		if res.Error != nil {
			err = fmt.Errorf(res.Error.Description)
		}
	}

	return err
}
