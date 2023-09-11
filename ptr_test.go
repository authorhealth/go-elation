package elation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	assert := assert.New(t)

	strVal := "foo"
	assert.Equal(&strVal, Ptr(strVal))

	boolVal := true
	assert.Equal(&boolVal, Ptr(boolVal))

	intVal := 2
	assert.Equal(&intVal, Ptr(intVal))
}
