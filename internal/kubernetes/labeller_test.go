package kubernetes

import (
	// standard packages
	"context"
	"testing"

	// external packages
	"github.com/stretchr/testify/assert"

	// internal packages
	testInternal "volume-cleaner/internal/testing"
)

func TestLabelFunctions(t *testing.T) {

	t.Run("test add, remove, edit labels", func(t *testing.T) {
		// create fake client
		kube := testInternal.NewFakeClient()

		labels := map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}
		if namespaceErr := kube.CreateNamespace(context.TODO(), "test", labels); namespaceErr != nil {
			t.Fatalf("Error injecting namespace add: %v", namespaceErr)
		}

		if _, pvcErr := kube.CreatePersistentVolumeClaim(context.TODO(), "pvc1", "test"); pvcErr != nil {
			t.Fatalf("Error injecting pvc add: %v", pvcErr)
		}

		// test adding new label
		SetPvcLabel(kube, "volume-cleaner/unattached-time", "foo", "test", "pvc1")

		assert.Equal(t, PvcList(kube, "test")[0].Labels["volume-cleaner/unattached-time"], "foo")

		// test changing existing label
		SetPvcLabel(kube, "volume-cleaner/unattached-time", "bar", "test", "pvc1")

		assert.Equal(t, PvcList(kube, "test")[0].Labels["volume-cleaner/unattached-time"], "bar")

		// test removing label
		RemovePvcLabel(kube, "volume-cleaner/unattached-time", "test", "pvc1")

		_, ok := PvcList(kube, "test")[0].Labels["volume-cleaner/unattached-time"]

		assert.Equal(t, ok, false)
	})
}
