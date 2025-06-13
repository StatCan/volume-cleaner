package kubernetes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestLabelFunctions(t *testing.T) {

	t.Run("test add, remove, edit labels", func(t *testing.T) {
		// create fake client
		client := fake.NewClientset()

		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test",
			Labels: map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}}}
		_, namespaceErr := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
		if namespaceErr != nil {
			t.Fatalf("Error injecting namespace add: %v", namespaceErr)
		}

		pvc := &corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: "pvc1", Namespace: "test"}}
		_, pvcErr := client.CoreV1().PersistentVolumeClaims("test").Create(context.TODO(), pvc, metav1.CreateOptions{})
		if pvcErr != nil {
			t.Fatalf("Error injecting pvc add: %v", pvcErr)
		}

		// test adding new label
		SetPvcLabel(client, "volume-cleaner/unattached-time", "foo", "test", "pvc1")

		assert.Equal(t, PvcList(client, "test")[0].Labels["volume-cleaner/unattached-time"], "foo")

		// test changing existing label
		SetPvcLabel(client, "volume-cleaner/unattached-time", "bar", "test", "pvc1")

		assert.Equal(t, PvcList(client, "test")[0].Labels["volume-cleaner/unattached-time"], "bar")

		// test removing label
		RemovePvcLabel(client, "volume-cleaner/unattached-time", "test", "pvc1")

		_, ok := PvcList(client, "test")[0].Labels["volume-cleaner/unattached-time"]

		assert.Equal(t, ok, false)
	})
}
