package utils

import (
	// External Imports
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {

	t.Run("successful", func(t *testing.T) {
		assert.Equal(t, true, true)

	})
}
