package elation

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeWithOptionalZone_UnmarshalJSON(t *testing.T) {
	laLoc, err := time.LoadLocation("America/Los_Angeles")
	assert.NoError(t, err)

	testCases := map[string]struct {
		jsonValue     string
		expectedError string
		expectedValue TimeWithOptionalZone
	}{
		"empty JSON value": {
			jsonValue:     "",
			expectedError: "unexpected end of JSON input",
		},
		"null JSON value": {
			jsonValue:     "null",
			expectedValue: TimeWithOptionalZone{},
		},
		"numeric JSON value": {
			jsonValue:     "12345",
			expectedError: "parsing raw time as JSON string",
		},
		"empty JSON string": {
			jsonValue:     `""`,
			expectedError: "parsing time",
		},
		"RFC3339 timestamp": {
			jsonValue:     `"2006-01-02T15:04:05+07:00"`,
			expectedValue: TimeWithOptionalZone(time.Date(2006, 1, 2, 8, 4, 5, 0, time.UTC)),
		},
		"ISO8601 timestamp without timezone": {
			jsonValue:     `"2006-01-02T15:04:05"`,
			expectedValue: TimeWithOptionalZone(time.Date(2006, 1, 2, 15, 4, 5, 0, laLoc)),
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			var value TimeWithOptionalZone
			err := json.Unmarshal([]byte(testCase.jsonValue), &value)
			if testCase.expectedError != "" {
				assert.ErrorContains(err, testCase.expectedError)
			} else {
				assert.NoError(err)
				assert.WithinDuration(testCase.expectedValue.Time(), value.Time(), 0)
			}
		})
	}
}
