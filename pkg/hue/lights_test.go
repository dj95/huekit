package hue

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func testServerLights() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.RequestURI, "/lights") {
			lightsHandler(w, r)

			return
		}

		lightHandler(w, r)
	}))
}

func lightsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
		return
	}

	body := `{
  "1": {
    "name": "TV Left"
  }
}`

	w.WriteHeader(200)
	w.Write([]byte(body))
}

func lightHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
		return
	}

	body := `{
  "state": {
    "on": true,
    "bri": 202,
    "hue": 13122,
    "sat": 211,
    "xy": [
      0.5119,
      0.4147
    ],
    "ct": 467,
    "alert": "none",
    "effect": "none",
    "colormode": "ct",
    "reachable": true
  },
  "type": "Extended color light",
  "name": "TV Left",
  "modelid": "LCT001",
  "swversion": "65003148",
  "pointsymbol": {
    "1": "none",
    "2": "none",
    "3": "none",
    "4": "none",
    "5": "none",
    "6": "none",
    "7": "none",
    "8": "none"
  }
}`

	w.WriteHeader(200)
	w.Write([]byte(body))
}

func TestBridge_Lights(t *testing.T) {
	mockServer := testServerLights()

	tests := []struct {
		description    string
		address        string
		username       string
		expectedError  bool
		expectedResult []*Light
	}{
		{
			description:   "success",
			address:       strings.TrimLeft(mockServer.URL, "htp:/"),
			username:      "success",
			expectedError: false,
			expectedResult: []*Light{
				{
					ID:              "1",
					Type:            "Extended color light",
					Name:            "TV Left",
					ModelID:         "LCT001",
					SoftwareVersion: "65003148",
					State: &State{
						On:               true,
						Brightness:       202,
						Hue:              13122,
						Saturation:       211,
						XY:               [2]float64{0.5119, 0.4147},
						ColorTemperature: 467,
						Alert:            "none",
						Effect:           "none",
						ColorMode:        "ct",
						Reachable:        true,
					},
					PointSymbol: map[string]string{
						"1": "none",
						"2": "none",
						"3": "none",
						"4": "none",
						"5": "none",
						"6": "none",
						"7": "none",
						"8": "none",
					},
				},
			},
		},
	}

	for _, test := range tests {
		bridge := &Bridge{
			address:  test.address,
			username: test.username,
		}

		result, err := bridge.Lights()

		assert.Equalf(t, test.expectedError, err != nil, test.description)
		assert.Nilf(t, deep.Equal(test.expectedResult, result), test.description)
	}
}
