package structure

import (
	// standard packages
	"testing"

	// external packages
	"github.com/stretchr/testify/assert"
)

func TestNewSet(t *testing.T) {

	t.Run("successful set creation", func(t *testing.T) {
		// check set was made
		actual := NewSet()

		expected := &Set{}
		expected.list = make(map[string]struct{})

		assert.Equal(t, expected, actual)

		// check that map was made
		assert.Equal(t, actual.list, make(map[string]struct{}))

	})
}

func TestSetMethods(t *testing.T) {

	t.Run("valid set operations", func(t *testing.T) {
		s := NewSet()

		// test add values
		s.Add("hello")
		s.Add("world")

		assert.Equal(t, s.Has("hello"), true)
		assert.Equal(t, s.Has("world"), true)
		assert.Equal(t, s.Has("go"), false)

		assert.Equal(t, s.Length(), 2)

		// test add duplicate value
		s.Add("hello")
		s.Add("go")

		assert.Equal(t, s.Has("hello"), true)
		assert.Equal(t, s.Has("go"), true)

		assert.Equal(t, s.Length(), 3)

		// test remove value
		s.Remove("hello")
		s.Remove("world")

		assert.Equal(t, s.Has("hello"), false)
		assert.Equal(t, s.Has("world"), false)
		assert.Equal(t, s.Has("go"), true)

		assert.Equal(t, s.Length(), 1)

		// test remove duplicate value
		s.Remove("hello")

		assert.Equal(t, s.Has("hello"), false)

		assert.Equal(t, s.Length(), 1)

		// test clear
		s.Clear()

		assert.Equal(t, s.Has("hello"), false)
		assert.Equal(t, s.Has("world"), false)
		assert.Equal(t, s.Has("go"), false)

		assert.Equal(t, s.Length(), 0)
	})
}

func TestSetDifference(t *testing.T) {
	type testCase struct {
		otherSet *Set
		expected *Set
	}

	t.Run("valid set difference", func(t *testing.T) {
		s := NewSet()

		s.Add("1")
		s.Add("2")
		s.Add("3")

		tests := []testCase{
			{
				NewSet(),
				func() *Set {
					s1 := NewSet()
					s1.Add("1")
					s1.Add("2")
					s1.Add("3")
					return s1
				}(),
			},
			{
				func() *Set {
					s2 := NewSet()
					s2.Add("1")
					return s2
				}(),
				func() *Set {
					s3 := NewSet()
					s3.Add("2")
					s3.Add("3")
					return s3
				}(),
			},
			{
				func() *Set {
					s4 := NewSet()
					s4.Add("1")
					s4.Add("2")
					s4.Add("3")
					return s4
				}(),
				NewSet(),
			},
		}

		for _, test := range tests {
			actual := s.Difference(test.otherSet)
			assert.Equal(t, test.expected, actual)
		}

		assert.Equal(t, NewSet().Difference(func() *Set {
			s := NewSet()
			s.Add("1")
			s.Add("2")
			s.Add("3")
			return s
		}()), NewSet())

	})
}
