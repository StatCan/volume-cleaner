package utils

import (
	// standard packages
	"testing"

	// external packages
	"github.com/stretchr/testify/assert"
)

// TestParseNotifTimes tests the ParseNotifTimes function with various inputs.
func TestParseNotifTimes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "single digit",
			input:    "1",
			expected: []int{1},
		},
		{
			name:     "multiple digits with space",
			input:    "1, 2",
			expected: []int{2, 1},
		},
		{
			name:     "multiple digits without space",
			input:    "1,2",
			expected: []int{2, 1},
		},
		{
			name:     "multiple digits",
			input:    "1,2,3",
			expected: []int{3, 2, 1},
		},
		{
			name:     "multiple digits with varied spacing",
			input:    "	1,      2   , 3 ",
			expected: []int{3, 2, 1},
		},
		{
			name:     "unsorted digits, should be sorted",
			input:    "3, 2, 1",
			expected: []int{3, 2, 1},
		},
		{
			name:     "empty string",
			input:    "",
			expected: make([]int, 0), // Assuming empty string results in an empty slice
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseNotifTimes(tt.input)
			assert.Equal(t, tt.expected, actual, "for input: %q", tt.input)
		})
	}
}

// TestParseGracePeriod tests the ParseGracePeriod function.
func TestParseGracePeriod(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "valid grace period",
			input:    "30",
			expected: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseGracePeriod(tt.input)
			assert.Equal(t, tt.expected, actual, "for input: %q", tt.input)
		})
	}
}
