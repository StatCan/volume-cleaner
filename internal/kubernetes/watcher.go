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
func WatchSts(ctx context.Context, kube *kubernetes.Clientset) {
	// leaving namespace as anray-liu for now until more rigorous testing is done
	watcher, err := kube.AppsV1().StatefulSets("anray-liu").Watch(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Fatalf("Error creating a watcher for statefulsets: %v", err)
	}

	log.Print("Watching for statefulset events...")

	// Watch Loop
	events := watcher.ResultChan()

	for {
		select {
		case <-ctx.Done():
			return

		case event := <-events:
			sts, ok := event.Object.(*appsv1.StatefulSet)

			// Skip this event if it can't be parsed into a sts
			if !ok {
				continue
			}

			// stsPvcs := PvcListBySts(kube, sts)

			switch event.Type {
			case watch.Added:
				log.Printf("sts added: %s", sts.Name)
				for _, vol := range sts.Spec.Template.Spec.Volumes {
					log.Printf("removing label")

					RemovePvcLabel(kube, "volume-cleaner/unattached-time", sts.Namespace, vol.PersistentVolumeClaim.ClaimName)
				}
			case watch.Deleted:
				log.Printf("sts deleted: %s", sts.Name)
				for _, vol := range sts.Spec.Template.Spec.Volumes {
					log.Printf("adding label")

					SetPvcLabel(kube, "volume-cleaner/unattached-time", time.Now().Format("2006-01-02_15-04-05Z"), sts.Namespace, vol.PersistentVolumeClaim.ClaimName)
				}
			}
		}
	}

}
