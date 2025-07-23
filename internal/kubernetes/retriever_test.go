package kubernetes

import (
	// standard packages
	"context"
	"testing"

	// external packages
	"github.com/stretchr/testify/assert"

	// internal packages
	"volume-cleaner/internal/structure"
	testInternal "volume-cleaner/internal/tests"
)

func TestNsList(t *testing.T) {

	t.Run("successful ns listing", func(t *testing.T) {
		// create fake client
		kube := testInternal.NewFakeClient()

		// set information for fake namespaces to be injected into client
		labels := map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}

		names := []string{"ns1", "ns2", "ns3", "ns4"}

		// inject fake namespaces
		for _, name := range names {
			if namespaceErr := kube.CreateNamespace(context.TODO(), name, labels); namespaceErr != nil {
				t.Fatalf("Error injecting namespace add: %v", namespaceErr)
			}
		}

		list := NsList(kube)

		// check right length
		assert.Equal(t, len(list), len(names))

		// check that each namespace is found
		for i, ns := range list {
			assert.Equal(t, ns.Name, names[i])
		}

	})
}

func TestStsList(t *testing.T) {

	t.Run("successful sts listing", func(t *testing.T) {
		// create fake client
		kube := testInternal.NewFakeClient()

		labels := map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}
		if namespaceErr := kube.CreateNamespace(context.TODO(), "test", labels); namespaceErr != nil {
			t.Fatalf("Error injecting namespace add: %v", namespaceErr)
		}

		names := []string{"sts1", "sts2", "sts3"}

		// inject fake stateful sets
		for _, name := range names {
			if stsErr := kube.CreateStatefulSet(context.TODO(), name, "test"); stsErr != nil {
				t.Fatalf("Error injecting sts add: %v", stsErr)
			}
		}

		list := StsList(kube, "test")

		// check right length
		assert.Equal(t, len(list), len(names))

		// check that each sts is found
		for i, sts := range list {
			assert.Equal(t, sts.Name, names[i])
		}

	})
}

func TestPvcList(t *testing.T) {

	t.Run("successful pvc listing", func(t *testing.T) {
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

		list := PvcList(kube, "test")

		// check right length
		assert.Equal(t, len(list), len(names))

		// check that each pvc is found
		for i, pvc := range list {
			assert.Equal(t, pvc.Name, names[i])
		}

	})
}

func TestFindUnattachedPVCs(t *testing.T) {

	t.Run("successfully find unattached pvcs", func(t *testing.T) {
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

		assert.Equal(t, len(FindUnattachedPVCs(kube, structure.ControllerConfig{})), 2)

		// mock a stateful set attached to a pvc1
		if stsErr := kube.CreateStatefulSetWithPvc(context.TODO(), "sts1", "test", "pvc1"); stsErr != nil {
			t.Fatalf("Error injecting sts add: %v", stsErr)
		}

		assert.Equal(t, len(FindUnattachedPVCs(kube, structure.ControllerConfig{})), 1)

		// mock a sts with no vols
		if err := kube.CreateStatefulSet(context.TODO(), "sts-no-volumes", "test"); err != nil {
			t.Fatalf("error creating sts: %v", err)
		}

		// no new attachements, expected unattached should still be 1
		assert.Equal(t, len(FindUnattachedPVCs(kube, structure.ControllerConfig{})), 1)

	})
}
