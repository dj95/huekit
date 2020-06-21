package hue

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Light Represents a light/plug at the hue bridge
type Light struct {
	Type            string            `json:"type"`
	Name            string            `json:"name"`
	ModelID         string            `json:"modelid"`
	SoftwareVersion string            `json:"swversion"`
	State           *State            `json:"state"`
	PointSymbol     map[string]string `json:"pointsymbol"`
}

// State Represents the state of a light
type State struct {
	On               bool       `json:"on"`
	Brightness       int        `json:"bri"`
	Hue              int        `json:"hue"`
	Saturation       int        `json:"sat"`
	XY               [2]float64 `json:"xy"`
	ColorTemperature int        `json:"ct"`
	Alert            string     `json:"alert"`
	Effect           string     `json:"effect"`
	ColorMode        string     `json:"colormode"`
	Reachable        bool       `json:"reachable"`
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
	bodyBytes, err := ioutil.ReadAll(res.Body)

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
	bodyBytes, err := ioutil.ReadAll(res.Body)

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

	return &light, nil
}
