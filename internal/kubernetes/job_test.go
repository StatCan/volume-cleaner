package kubernetes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsStale(t *testing.T) {

	t.Run("test successful determination of stale pvcs", func(t *testing.T) {
		format := "2006-01-02_15-04-05Z"

		assert.Equal(t, IsStale(time.Now().Add(-time.Hour*24*180).Format(format), format, 180), false)
		assert.Equal(t, IsStale(time.Now().Add(-time.Hour*24*181).Format(format), format, 180), true)
		assert.Equal(t, IsStale(time.Now().Add(-time.Hour*24*180-time.Hour*23).Format(format), format, 180), false)
		assert.Equal(t, IsStale(time.Now().Add(-time.Hour*24*180-time.Hour*23-time.Minute*59-time.Second*59).Format(format), format, 180), false)
		assert.Equal(t, IsStale(time.Now().Add(-time.Hour*24*1000).Format(format), format, 180), true)
		assert.Equal(t, IsStale(time.Now().Format(format), format, 180), false)
		assert.Equal(t, IsStale(time.Now().Format(format), format, 0), false)
		assert.Equal(t, IsStale(time.Now().Add(-time.Second).Format(format), format, 0), false)
		assert.Equal(t, IsStale(time.Now().Add(-time.Hour*24).Format(format), format, 0), true)

	})
}
