package kubernetes

import (
	// External Imports
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNsList(t *testing.T) {

	t.Run("successful ns listing", func(t *testing.T) {
		// create fake client
		client := fake.NewClientset()

		// set information for fake namespaces to be injected into client
		labels := make(map[string]string)
		labels["app.kubernetes.io/part-of"] = "kubeflow-profile"

		names := []string{"ns1", "ns2", "ns3", "ns4"}

		// inject fake namespaces
		for _, name := range names {
			ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name, Labels: labels}}
			_, err := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Error injecting namespace add: %v", err)
			}
		}

		list := NsList(client)

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
		client := fake.NewClientset()

		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test",
			Labels: map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}}}
		_, err := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Error injecting namespace add: %v", err)
		}

		names := []string{"sts1", "sts2", "sts3"}

		// inject fake stateful sets
		for _, name := range names {
			sts := &appv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "test"}}
			_, err := client.AppsV1().StatefulSets("test").Create(context.TODO(), sts, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Error injecting sts add: %v", err)
			}
		}

		list := StsList(client, "test")

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
		client := fake.NewClientset()

		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test",
			Labels: map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}}}
		_, err := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Error injecting namespace add: %v", err)
		}

		names := []string{"pvc1", "pvc2"}

		// inject fake pvcs
		for _, name := range names {
			pvc := &corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "test"}}
			_, err := client.CoreV1().PersistentVolumeClaims("test").Create(context.TODO(), pvc, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Error injecting pvc add: %v", err)
			}
		}

		list := PvcList(client, "test")

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
		client := fake.NewClientset()

		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test",
			Labels: map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}}}
		_, err := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Error injecting namespace add: %v", err)
		}

		names := []string{"pvc1", "pvc2"}

		// inject fake pvcs
		for _, name := range names {
			pvc := &corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "test"}}
			_, err := client.CoreV1().PersistentVolumeClaims("test").Create(context.TODO(), pvc, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Error injecting pvc add: %v", err)
			}
		}

		assert.Equal(t, len(FindUnattachedPVCs(client)), 2)

		// mock a stateful set attached to a pvc1

		sts := &appv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sts1",
				Namespace: "test",
			},
			Spec: appv1.StatefulSetSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Volumes: []corev1.Volume{
							{
								Name: "pvc1",
								VolumeSource: corev1.VolumeSource{
									PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
										ClaimName: "pvc1",
									},
								},
							},
						},
					},
				},
			},
		}

		_, err = client.AppsV1().StatefulSets("test").Create(context.TODO(), sts, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Error injecting sts add: %v", err)

		}

		assert.Equal(t, len(FindUnattachedPVCs(client)), 1)

	})
}
