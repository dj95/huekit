package hue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsModelIDFromHue(t *testing.T) {
	tests := []struct {
		description    string
		modelID        string
		expectedResult bool
	}{
		{
			description:    "model id from hue",
			modelID:        "LTW012",
			expectedResult: true,
		},
		{
			description:    "model id not from hue",
			modelID:        "Plug",
			expectedResult: false,
		},
	}

	for _, test := range tests {
		result := ModelIDIsFromHue(test.modelID)

		assert.Equalf(t, test.expectedResult, result, test.description)
	}
}
