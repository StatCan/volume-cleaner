package kubernetes

import (
	// External Imports
	"context"
	"log"
	"time"

	structInternal "volume-cleaner/internal/structure"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// Watches for when statefulsets are created or deleted across all namespaces
func WatchSts(ctx context.Context, kube kubernetes.Interface, cfg structInternal.Config) {
	// leaving namespace as anray-liu for now until more rigorous testing is done
	// reminder to not hard code namspace after unit tests are done

	watcher, err := kube.AppsV1().StatefulSets(cfg.Namespace).Watch(ctx, metav1.ListOptions{})

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

func InitialScan(kube kubernetes.Interface, cfg structInternal.Config) {
	log.Print("Checking for unattached PVCs...")
	for _, pvc := range FindUnattachedPVCs(kube) {
		_, ok := pvc.Labels[cfg.Label]
		if !ok {
			log.Print("PVC labelled.")
			SetPvcLabel(kube, cfg.Label, time.Now().Format(cfg.TimeFormat), pvc.Namespace, pvc.Name)
		} else {
			log.Print("PVC is attached. Skipping.")
		}
	}
}
