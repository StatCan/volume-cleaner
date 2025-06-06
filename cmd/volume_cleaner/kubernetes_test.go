package main

import (
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
			ns := &appv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "test"}}
			_, err := client.AppsV1().StatefulSets("test").Create(context.TODO(), ns, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Error injecting namespace add: %v", err)
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
