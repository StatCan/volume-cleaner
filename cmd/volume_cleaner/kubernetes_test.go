package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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

		names := []string{"n1", "n2", "n3", "n4"}

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
		assert.Equal(t, len(list), 4)

		// check that each namespace is found
		for i, ns := range list {
			assert.Equal(t, ns.Name, names[i])
		}

	})
}
