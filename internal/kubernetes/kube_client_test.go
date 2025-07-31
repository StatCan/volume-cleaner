package kubernetes

import (
	// standard packages
	"reflect"
	"testing"

	// external packages
	"github.com/stretchr/testify/assert"
)

func TestInitClient(t *testing.T) {

	t.Run("test successful creation of kube client", func(t *testing.T) {
		kube, err := InitKubeClient()
		assert.NotEqual(t, nil, err) // will throw error when not running on a cluster
		assert.Equal(t, "*kubernetes.Clientset", reflect.TypeOf(kube).String())
	})
}
