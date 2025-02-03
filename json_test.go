package elation

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonNullJSONArray_MarshalJSON(t *testing.T) {
	testCases := map[string]struct {
		value         NonNullJSONArray[string]
		expectedBytes []byte
	}{
		"nil value": {
			value:         nil,
			expectedBytes: []byte("[]"),
		},
		"empty array": {
			value:         NonNullJSONArray[string]{},
			expectedBytes: []byte("[]"),
		},
		"non-empty array": {
			value:         NonNullJSONArray[string]{"alvin", "simon", "theodore"},
			expectedBytes: []byte(`["alvin","simon","theodore"]`),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			bytes, err := testCase.value.MarshalJSON()
			assert.NoError(err)

			assert.Equal(testCase.expectedBytes, bytes)
		})
	}
}

func TestNonNullJSONArray_UnmarshalJSON(t *testing.T) {
	testCases := map[string]struct {
		bytes         []byte
		expectedValue NonNullJSONArray[string]
	}{
		"empty array": {
			bytes:         []byte("[]"),
			expectedValue: nil,
		},
		"null": {
			bytes:         []byte("null"),
			expectedValue: nil,
		},
		"non-empty array": {
			bytes:         []byte(`["foo", "bar"]`),
			expectedValue: NonNullJSONArray[string]{"foo", "bar"},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			var val NonNullJSONArray[string]
			err := json.Unmarshal(testCase.bytes, &val)
			assert.NoError(err)

			assert.Equal(testCase.expectedValue, val)
		})
	}
}
