package utils

import (
	// standard packages
	"context"
	"testing"

	// external packages
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// TestNewFakeClient verifies the creation of a K8S client
func TestNewFakeClient(t *testing.T) {
	f := NewFakeClient()
	assert.NotNil(t, f, "NewFakeClient should return a non-nil FakeClient")
	// Ensure embedded Interface is a fake clientset
	_, ok := f.Interface.(*fake.Clientset)
	assert.True(t, ok, "Interface should be a *fake.Clientset")
}

// TestCreateNamespace verifies the creation of a namespace in the fake K8S client
func TestCreateNamespace(t *testing.T) {
	ctx := context.TODO()
	f := NewFakeClient()
	name := "test-ns"
	labels := map[string]string{"env": "test"}
	// create namespace
	err := f.CreateNamespace(ctx, name, labels)
	assert.NoError(t, err, "expected no error creating namespace")

	// verify namespace exists with correct labels
	ns, err := f.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Equal(t, name, ns.Name)
	assert.Equal(t, labels, ns.Labels)
}

// TestCreatePersistentVolumeClaim verifies the creation of a Persistent Volume Claim within a namespace
func TestCreatePersistentVolumeClaim(t *testing.T) {
	ctx := context.TODO()
	f := NewFakeClient()
	ns := "pvc-ns"
	// need namespace before PVC
	err := f.CreateNamespace(ctx, ns, nil)
	assert.NoError(t, err)

	pvcName := "test-pvc"
	pvc, err := f.CreatePersistentVolumeClaim(ctx, pvcName, ns)
	assert.NoError(t, err, "expected no error creating PVC")
	assert.Equal(t, pvcName, pvc.Name)
	assert.Equal(t, ns, pvc.Namespace)

	// verify via client
	got, err := f.CoreV1().PersistentVolumeClaims(ns).Get(ctx, pvcName, metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Equal(t, pvcName, got.Name)
}

// TestCreateStatefulSet verifies the creation of a statefulset within a namespace
func TestCreateStatefulSet(t *testing.T) {
	ctx := context.TODO()
	f := NewFakeClient()
	ns := "sts-ns"
	err := f.CreateNamespace(ctx, ns, nil)
	assert.NoError(t, err)

	stsName := "test-sts"
	err = f.CreateStatefulSet(ctx, stsName, ns)
	assert.NoError(t, err, "expected no error creating StatefulSet")

	got, err := f.AppsV1().StatefulSets(ns).Get(ctx, stsName, metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Equal(t, stsName, got.Name)
}

// TestCreateStatefulSetWithPvc verifies the creation of a statefulset within a namespace along with a bounded Persistent Volume Claim
func TestCreateStatefulSetWithPvc(t *testing.T) {
	ctx := context.TODO()
	f := NewFakeClient()
	ns := "sts-pvc-ns"
	// create namespace and PVC
	err := f.CreateNamespace(ctx, ns, nil)
	assert.NoError(t, err)
	pvcName := "data-pvc"
	_, err = f.CreatePersistentVolumeClaim(ctx, pvcName, ns)
	assert.NoError(t, err)

	stsName := "sts-with-pvc"
	err = f.CreateStatefulSetWithPvc(ctx, stsName, ns, pvcName)
	assert.NoError(t, err, "expected no error creating StatefulSet with PVC")

	got, err := f.AppsV1().StatefulSets(ns).Get(ctx, stsName, metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Equal(t, stsName, got.Name)
	// verify volume reference
	vols := got.Spec.Template.Spec.Volumes
	assert.Len(t, vols, 1)
	assert.Equal(t, pvcName, vols[0].Name)
	assert.NotNil(t, vols[0].PersistentVolumeClaim)
	assert.Equal(t, pvcName, vols[0].PersistentVolumeClaim.ClaimName)
}

// TestDeleteStatefulSet verifies the deletion of a statefulset within a namespace
func TestDeleteStatefulSet(t *testing.T) {
	ctx := context.TODO()
	f := NewFakeClient()

	ns := "del-sts-ns"
	err := f.CreateNamespace(ctx, ns, nil)
	assert.NoError(t, err)

	stsName := "to-delete-sts"
	err = f.CreateStatefulSet(ctx, stsName, ns)
	assert.NoError(t, err)

	// ensure exists
	_, err = f.AppsV1().StatefulSets(ns).Get(ctx, stsName, metav1.GetOptions{})
	assert.NoError(t, err)

	// delete
	err = f.DeleteStatefulSet(ctx, stsName, ns)
	assert.NoError(t, err, "expected no error deleting StatefulSet")

	// verify deletion
	_, err = f.AppsV1().StatefulSets(ns).Get(ctx, stsName, metav1.GetOptions{})
	assert.Error(t, err, "expected error getting deleted StatefulSet")
}
