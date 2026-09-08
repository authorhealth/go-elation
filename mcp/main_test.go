package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseServerOptions(t *testing.T) {
	t.Run("defaults to safe tools only", func(t *testing.T) {
		assert := assert.New(t)
		opts, err := parseServerOptions(nil)
		assert.NoError(err)
		assert.False(opts.allowUnsafeTools)
	})

	t.Run("enables unsafe tools flag", func(t *testing.T) {
		assert := assert.New(t)
		opts, err := parseServerOptions([]string{"--allow-unsafe-tools"})
		assert.NoError(err)
		assert.True(opts.allowUnsafeTools)
	})
}
