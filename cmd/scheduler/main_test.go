package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNotifTimes(t *testing.T) {

	t.Run("test successful parsing of notification times", func(t *testing.T) {
		assert.Equal(t, ParseNotifTimes("1"), []int{1})
		assert.Equal(t, ParseNotifTimes("1, 2"), []int{1, 2})
		assert.Equal(t, ParseNotifTimes("1,2"), []int{1, 2})
		assert.Equal(t, ParseNotifTimes("1,2,3"), []int{1, 2, 3})
		assert.Equal(t, ParseNotifTimes("	1,      2   , 3 "), []int{1, 2, 3})
		assert.Equal(t, ParseNotifTimes("3, 2, 1"), []int{1, 2, 3})
	})
}

func TestParseGracePeriod(t *testing.T) {

	t.Run("test successful parsing of grace period", func(t *testing.T) {

		assert.Equal(t, ParseGracePeriod("30"), 30)

	})

}
