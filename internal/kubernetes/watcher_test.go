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
	testInternal "volume-cleaner/internal/utils"
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
			TimeLabel:  "volume-cleaner/unattached-time",
			NotifLabel: "volume-cleaner/notification-count",
			TimeFormat: "2006-01-02_15-04-05Z",
		}

		go WatchSts(ctx, kube, cfg)

		time.Sleep(2 * time.Second)

		// no pvc should have labels right now
		pvcs := PvcList(kube, "test")

		_, ok := pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[0].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/notification-count"]
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

		_, ok = pvcs[0].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/notification-count"]
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

		_, ok = pvcs[0].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, true)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, false)

		// add back sts1
		if stsErr := kube.CreateStatefulSetWithPvc(context.TODO(), "sts1", "test", "pvc1"); stsErr != nil {
			t.Fatalf("Error injecting sts add: %v", stsErr)
		}

		time.Sleep(2 * time.Second)

		pvcs = PvcList(kube, "test")

		_, ok = pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[0].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, false)

		ctx.Done()

	})
}

func TestWatcherStorageClassFilter(t *testing.T) {

	t.Run("successful skipping of unconfigured storage classes", func(t *testing.T) {
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
			Namespace:      "test",
			TimeLabel:      "volume-cleaner/unattached-time",
			NotifLabel:     "volume-cleaner/notification-count",
			TimeFormat:     "2006-01-02_15-04-05Z",
			StorageClasses: []string{"non-existent-storage-class"},
		}

		go WatchSts(ctx, kube, cfg)

		// mock a stateful set attached to a pvc1
		if stsErr := kube.CreateStatefulSetWithPvc(context.TODO(), "sts1", "test", "pvc1"); stsErr != nil {
			t.Fatalf("Error injecting sts add: %v", stsErr)
		}

		// delete sts
		if eventErr := kube.DeleteStatefulSet(context.TODO(), "sts1", "test"); eventErr != nil {
			t.Fatalf("Error injecting event add: %v", eventErr)
		}

		time.Sleep(2 * time.Second)

		// should not have new labels

		pvcs := PvcList(kube, "test")

		_, ok := pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[0].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/notification-count"]
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

		_, ok = pvcs[0].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, false)

		cfg := structInternal.ControllerConfig{
			Namespace:  "test",
			TimeLabel:  "volume-cleaner/unattached-time",
			NotifLabel: "volume-cleaner/notification-count",
			TimeFormat: "2006-01-02_15-04-05Z",
		}

		InitialScan(kube, cfg)

		time.Sleep(2 * time.Second)

		// should have new labels

		pvcs = PvcList(kube, "test")

		_, ok = pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, true)

		_, ok = pvcs[0].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, true)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, true)

		_, ok = pvcs[1].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, true)

	})
}

func TestResetLabels(t *testing.T) {

	t.Run("successful resetting of labels on controller startup", func(t *testing.T) {
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

		cfg := structInternal.ControllerConfig{
			Namespace:  "test",
			TimeLabel:  "volume-cleaner/unattached-time",
			NotifLabel: "volume-cleaner/notification-count",
			TimeFormat: "2006-01-02_15-04-05Z",
		}

		InitialScan(kube, cfg)

		// pvcs should be labelled
		pvcs := PvcList(kube, "test")

		_, ok := pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, true)

		_, ok = pvcs[0].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, true)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, true)

		_, ok = pvcs[1].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, true)

		ResetLabels(kube, cfg)

		time.Sleep(2 * time.Second)

		// should have all labels removed

		pvcs = PvcList(kube, "test")

		_, ok = pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[0].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/notification-count"]
		assert.Equal(t, ok, false)

	})
}

func TestIgnoreStorageClass(t *testing.T) {
	getPtr := func(str string) *string {
		variable := str
		return &variable
	}
	tests := []struct {
		name           string
		input          *string
		storageClasses []string
		expected       bool
	}{
		{
			name:           "single value",
			input:          getPtr("standard"),
			storageClasses: []string{"standard"},
			expected:       false,
		},
		{
			name:           "two values",
			input:          getPtr("standard"),
			storageClasses: []string{"standard", "default"},
			expected:       false,
		},
		{
			name:           "accept all with standard",
			input:          getPtr("standard"),
			storageClasses: []string{},
			expected:       false,
		},
		{
			name:           "accept all with empty",
			input:          nil,
			storageClasses: []string{},
			expected:       false,
		},
		{
			name:           "accept all with default",
			input:          getPtr("default"),
			storageClasses: []string{},
			expected:       false,
		},
		{
			name:           "reject single value",
			input:          getPtr("standard"),
			storageClasses: []string{"default"},
			expected:       true,
		},
		{
			name:           "reject two value",
			input:          getPtr("test"),
			storageClasses: []string{"standard", "default"},
			expected:       true,
		},
		{
			name:           "empty value",
			input:          getPtr(""),
			storageClasses: []string{""},
			expected:       false,
		},
		{
			name:           "accept nil value",
			input:          nil,
			storageClasses: []string{""},
			expected:       false,
		},
		{
			name:           "accept multiple values",
			input:          nil,
			storageClasses: []string{"default", "standard", ""},
			expected:       false,
		},
		{
			name:           "reject nil value",
			input:          nil,
			storageClasses: []string{"default"},
			expected:       true,
		},
		{
			name:           "reject whitespace",
			input:          nil,
			storageClasses: []string{" "},
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := IgnoreStorageClass(tt.input, tt.storageClasses)
			assert.Equal(t, tt.expected, actual, "for input: %q", tt.input)
		})
	}
}
