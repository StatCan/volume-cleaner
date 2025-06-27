package kubernetes

import (
	// standard packages
	"context"
	"log"
	"time"

	// external packages
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	// internal packages
	structInternal "volume-cleaner/internal/structure"
)

// Watches for when statefulsets are created or deleted across all namespaces

func WatchSts(ctx context.Context, kube kubernetes.Interface, cfg structInternal.ControllerConfig) {
	watcher, err := kube.AppsV1().StatefulSets(cfg.Namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error creating a watcher for statefulsets: %s", err)
	}

	log.Print("Watching for statefulset events...")

	events := watcher.ResultChan()

	for {
		select {

		// context used to kill loop
		case <-ctx.Done():
			return

		// sts was added or deleted
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

					RemovePvcLabel(kube, cfg.Label, sts.Namespace, vol.PersistentVolumeClaim.ClaimName)
				}
			case watch.Deleted:
				log.Printf("sts deleted: %s", sts.Name)

				for _, vol := range sts.Spec.Template.Spec.Volumes {
					log.Printf("adding label")

					SetPvcLabel(kube, cfg.Label, time.Now().Format(cfg.TimeFormat), sts.Namespace, vol.PersistentVolumeClaim.ClaimName)
				}
			}
		}
	}

}

func InitialScan(kube kubernetes.Interface, cfg structInternal.ControllerConfig) {
	log.Print("Checking for unattached PVCs...")
	for _, pvc := range FindUnattachedPVCs(kube) {
		_, ok := pvc.Labels[cfg.Label]
		if !ok {
			SetPvcLabel(kube, cfg.Label, time.Now().Format(cfg.TimeFormat), pvc.Namespace, pvc.Name)
		} else {
			log.Print("PVC already has label. Skipping.")
		}
	}
	log.Print("Initial scan complete")
}
