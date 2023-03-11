package hue

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testServerAuth() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`[]`))
			return
		}

		reqBytes, err := io.ReadAll(r.Body)

		if err != nil {
			fmt.Printf("[mock] cannot read body\n")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`[]`))
			return
		}

		var reqBody authRequest
		err = json.Unmarshal(reqBytes, &reqBody)

		if err != nil {
			fmt.Printf("[mock] cannot unmarshal body\n")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`[]`))
			return
		}

		body := `[{"success": {"username":"success"}}]`
		if strings.HasSuffix(reqBody.DeviceType, "#evil") {
			body = `[
  {
    "error": {
      "type": 7,
      "address": "/username",
      "description": "invalid value, burges, for parameter, username"
    }
  },
  {
    "error": {
      "type": 2,
      "address": "/",
      "description": "body contains invalid json"
    }
  }
]`
		}

		if strings.HasSuffix(reqBody.DeviceType, "#notpressed") {
			body = `[{"error": {"type": 101,"address": "","description": "link button not pressed"}}]`
		}

		w.WriteHeader(200)
		w.Write([]byte(body))
	}))
}

func TestPerformAuthRequest(t *testing.T) {
	mockServer := testServerAuth()

	tests := []struct {
		description    string
		address        string
		username       string
		expectedError  bool
		expectedResult string
	}{
		{
			description:    "success",
			address:        strings.TrimLeft(mockServer.URL, "htp:/"),
			username:       "success",
			expectedError:  false,
			expectedResult: "success",
		},
		{
			description:    "link button not pressed",
			address:        strings.TrimLeft(mockServer.URL, "htp:/"),
			username:       "notpressed",
			expectedError:  true,
			expectedResult: "",
		},
	}

	for _, test := range tests {
		result, err := performAuthRequest(
			test.address,
			test.username,
		)

		assert.Equalf(t, test.expectedError, err != nil, test.description)
		assert.Equalf(t, test.expectedResult, result, test.description)
	}
}
