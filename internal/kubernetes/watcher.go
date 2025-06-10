package kubernetes

import (
	// External Imports
	"context"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// Watches for when statefulsets are created or deleted across all namespaces
func WatchSts(kube *kubernetes.Clientset) {
	watcher, err := kube.AppsV1().StatefulSets("anray-liu").Watch(context.TODO(), metav1.ListOptions{})

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

		// stsPvcs := PvcListBySts(kube, sts)

		switch event.Type {
		case watch.Added:
			log.Printf("sts added: %s\n", sts.Name)
			for _, pvc := range PvcListBySts(kube, sts) {
				log.Printf("removing label from sts")
				RemovePvcLabel(kube, "volume-cleaner/unattached-time", "anray-liu", pvc.Name)
			}
		case watch.Deleted:
			log.Printf("sts deleted: %s\n", sts.Name)
			for _, vol := range sts.Spec.Template.Spec.Volumes {
				log.Printf("adding label to sts")
				SetPvcLabel(kube, "volume-cleaner/unattached-time", time.Now().String(), "anray-liu", vol.Name)
			}
		}

	}

}
