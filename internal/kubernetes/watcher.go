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

// Watches for when statefulsets are created or deleted

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
		// used during unit tests
		case <-ctx.Done():
			return

		// sts was added or deleted
		case event := <-events:
			sts, ok := event.Object.(*appsv1.StatefulSet)

			// Skip this event if it can't be parsed into a sts
			if !ok {
				log.Printf("%s could not be parsed into a STS", event)
				continue
			}

			switch event.Type {
			case watch.Added:
				log.Printf("sts added: %s", sts.Name)

				for _, vol := range sts.Spec.Template.Spec.Volumes {
					if vol.PersistentVolumeClaim == nil {
						continue
					}

					pvcObj, err := kube.CoreV1().PersistentVolumeClaims(sts.Namespace).Get(context.TODO(), vol.PersistentVolumeClaim.ClaimName, metav1.GetOptions{})
					if err != nil {
						log.Printf("Error finding PVC object: %s", err)
						continue
					}

					_, ok := pvcObj.Labels[cfg.TimeLabel]
					if ok {
						log.Printf("removing label %s", cfg.TimeLabel)
						RemovePvcLabel(kube, cfg.TimeLabel, sts.Namespace, vol.PersistentVolumeClaim.ClaimName)
					}

					_, ok = pvcObj.Labels[cfg.NotifLabel]
					if ok {
						log.Printf("removing label %s", cfg.NotifLabel)
						RemovePvcLabel(kube, cfg.NotifLabel, sts.Namespace, vol.PersistentVolumeClaim.ClaimName)
					}

				}
			case watch.Deleted:
				log.Printf("sts deleted: %s", sts.Name)

				for _, vol := range sts.Spec.Template.Spec.Volumes {
					if vol.PersistentVolumeClaim == nil {
						continue
					}

					log.Printf("adding labels")

					SetPvcLabel(kube, cfg.TimeLabel, time.Now().Format(cfg.TimeFormat), sts.Namespace, vol.PersistentVolumeClaim.ClaimName)
					SetPvcLabel(kube, cfg.NotifLabel, "0", sts.Namespace, vol.PersistentVolumeClaim.ClaimName)
				}
			}
		}
	}

}

func InitialScan(kube kubernetes.Interface, cfg structInternal.ControllerConfig) {
	log.Print("Checking for unattached PVCs...")

	// scan and repair labels for all unattached pvcs
	for _, pvc := range FindUnattachedPVCs(kube) {
		_, ok := pvc.Labels[cfg.TimeLabel]
		if !ok {
			log.Printf("adding missing label %s", cfg.TimeLabel)
			SetPvcLabel(kube, cfg.TimeLabel, time.Now().Format(cfg.TimeFormat), pvc.Namespace, pvc.Name)
		}
		_, ok = pvc.Labels[cfg.NotifLabel]
		if !ok {
			log.Printf("adding missing label %s", cfg.TimeLabel)
			SetPvcLabel(kube, cfg.NotifLabel, "0", pvc.Namespace, pvc.Name)
		}
	}

	log.Print("Initial scan complete")
}
