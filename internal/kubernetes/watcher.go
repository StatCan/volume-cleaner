package kubernetes

import (
	// External Imports
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// Watches for when statefulsets are created or deleted across all namespaces
func watchSts(kube *kubernetes.Clientset) {
	watcher, err := kube.AppsV1().StatefulSets(metav1.NamespaceAll).Watch(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Fatalf("Error creating a watcher for statefulsets: %v", err)
	}

	log.Print("Watching for statefulset events...")

	// Watch Loop
	events := watcher.ResultChan()
	for event := range events {
		sts, ok := event.Object.(*appsv1.StatefulSet)

		// Skip this event if it can't be parsed into a sts
		if !ok {
			continue
		}

		switch event.Type {
		case (watch.Added || watch.Deleted):
			// fmt.Printf("Pod added/deleted: %s\n", sts)
		}

	}

}
