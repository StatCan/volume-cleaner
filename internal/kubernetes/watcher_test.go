package kubernetes

import (
	// standard packages
	"context"
	"testing"
	"time"

	// external packages
	"github.com/stretchr/testify/assert"

	// internal packages
	structInternal "volume-cleaner/internal/structure"
	testInternal "volume-cleaner/internal/testing"
)

func TestWatcherLabelling(t *testing.T) {

	t.Run("successful labelling of attached and detached pvcs", func(t *testing.T) {
		// create fake client
		kube := testInternal.NewFakeClient()

		labels := map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}
		if namespaceErr := kube.CreateNamespace(context.TODO(), "test", labels); namespaceErr != nil {
			t.Fatalf("Error injecting namespace add: %v", namespaceErr)
		}

		names := []string{"pvc1", "pvc2"}

		// inject fake pvcs
		for _, name := range names {
			if _, pvcErr := kube.CreatePersistentVolumeClaim(context.TODO(), name, "test"); pvcErr != nil {
				t.Fatalf("Error injecting pvc add: %v", pvcErr)
			}
		}

		ctx := context.Background()

		cfg := structInternal.ControllerConfig{
			Namespace:  "test",
			Label:      "volume-cleaner/unattached-time",
			TimeFormat: "2006-01-02_15-04-05Z",
		}

		go WatchSts(ctx, kube, cfg)

		time.Sleep(2 * time.Second)

		// no pvc should have labels right now
		pvcs := PvcList(kube, "test")

		_, ok := pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		// mock a stateful set attached to a pvc1
		if stsErr := kube.CreateStatefulSetWithPvc(context.TODO(), "sts1", "test", "pvc1"); stsErr != nil {
			t.Fatalf("Error injecting sts add: %v", stsErr)
		}

		time.Sleep(2 * time.Second)

		// should be no change

		pvcs = PvcList(kube, "test")

		_, ok = pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		// delete sts
		if eventErr := kube.DeleteStatefulSet(context.TODO(), "sts1", "test"); eventErr != nil {
			t.Fatalf("Error injecting event add: %v", eventErr)
		}

		time.Sleep(2 * time.Second)

		// should have new labels

		pvcs = PvcList(kube, "test")

		_, ok = pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, true)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		ctx.Done()

	})
}

func TestInitialScan(t *testing.T) {

	t.Run("successful labelling of unatatched pvcs on controller startup", func(t *testing.T) {
		// create fake client
		kube := testInternal.NewFakeClient()

		labels := map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}
		if namespaceErr := kube.CreateNamespace(context.TODO(), "test", labels); namespaceErr != nil {
			t.Fatalf("Error injecting namespace add: %v", namespaceErr)
		}

		names := []string{"pvc1", "pvc2"}

		// inject fake pvcs
		for _, name := range names {
			if _, pvcErr := kube.CreatePersistentVolumeClaim(context.TODO(), name, "test"); pvcErr != nil {
				t.Fatalf("Error injecting pvc add: %v", pvcErr)
			}
		}

		// no pvc should have labels right now
		pvcs := PvcList(kube, "test")

		_, ok := pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		ctx := context.Background()

		cfg := structInternal.ControllerConfig{
			Namespace:  "test",
			Label:      "volume-cleaner/unattached-time",
			TimeFormat: "2006-01-02_15-04-05Z",
		}

		InitialScan(kube, cfg)

		time.Sleep(2 * time.Second)

		// should have new labels

		pvcs = PvcList(kube, "test")

		_, ok = pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, true)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, true)

		ctx.Done()

	})
}
