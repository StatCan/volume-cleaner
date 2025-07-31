package kubernetes

import (
	// standard packages
	"context"
	"log"
	"slices"
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
		log.Fatalf("[ERROR] Failed to create watcher for statefulsets: %s", err)
	}

	log.Print("[INFO] Watching for statefulset events...")

	// create a channel to capture sts events in the cluster
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
				continue
			}

			switch event.Type {

			case watch.Added:
				// sts added
				handleAdded(kube, cfg, sts)
			case watch.Deleted:
				// sts deleted
				handleDeleted(kube, cfg, sts)
			}
		}
	}

}

// scan performed on controller startup to find unattached pvcs and assign labels to them

func InitialScan(kube kubernetes.Interface, cfg structInternal.ControllerConfig) {
	log.Print("[INFO] Starting initial scan...")

	for _, pvc := range FindUnattachedPVCs(kube, cfg) {

		// add time stamp label if not found
		_, ok := pvc.Labels[cfg.TimeLabel]
		if !ok {
			log.Printf("[INFO] Adding missing label %s to %s", cfg.TimeLabel, pvc.Name)
			SetPvcLabel(kube, cfg.TimeLabel, time.Now().Format(cfg.TimeFormat), pvc.Namespace, pvc.Name)
		}

		// add notification count label if not found
		_, ok = pvc.Labels[cfg.NotifLabel]
		if !ok {
			log.Printf("[INFO] Adding missing label %s to %s", cfg.NotifLabel, pvc.Name)
			SetPvcLabel(kube, cfg.NotifLabel, "0", pvc.Namespace, pvc.Name)
		}
	}

	log.Print("[INFO] Initial scan complete.")
}

// scans all pvcs and removes all volume-cleaner related labels
func ResetLabels(kube kubernetes.Interface, cfg structInternal.ControllerConfig) {
	log.Print("Resetting labels...")

	for _, namespace := range NsList(kube) {
		for _, pvc := range PvcList(kube, namespace.Name) {
			_, ok := pvc.Labels[cfg.TimeLabel]
			if ok {
				RemovePvcLabel(kube, cfg.TimeLabel, namespace.Name, pvc.Name)
			}
			_, ok = pvc.Labels[cfg.NotifLabel]
			if ok {
				RemovePvcLabel(kube, cfg.NotifLabel, namespace.Name, pvc.Name)
			}

		}
	}
}

// triggered on sts creation event
// will remove labels from all associated pvcs

func handleAdded(kube kubernetes.Interface, cfg structInternal.ControllerConfig, sts *appsv1.StatefulSet) {
	log.Printf("[INFO] STS added: %s", sts.Name)

	for _, vol := range sts.Spec.Template.Spec.Volumes {
		if vol.PersistentVolumeClaim == nil {
			continue
		}

		// get pvc object from name
		pvcObj, err := kube.CoreV1().PersistentVolumeClaims(sts.Namespace).Get(context.TODO(), vol.PersistentVolumeClaim.ClaimName, metav1.GetOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to find PVC object %s: %s", vol.PersistentVolumeClaim.ClaimName, err)
			continue
		}
		log.Printf("[INFO] Found PVC object %s", pvcObj.Name)

		// ignore if storage class not in config
		if IgnoreStorageClass(pvcObj.Spec.StorageClassName, cfg.StorageClass) {
			continue
		}

		// remove labels if found
		_, ok := pvcObj.Labels[cfg.TimeLabel]
		if ok {
			log.Printf("[INFO] Removing label %s", cfg.TimeLabel)
			RemovePvcLabel(kube, cfg.TimeLabel, sts.Namespace, vol.PersistentVolumeClaim.ClaimName)
		}

		_, ok = pvcObj.Labels[cfg.NotifLabel]
		if ok {
			log.Printf("[INFO] Removing label %s", cfg.NotifLabel)
			RemovePvcLabel(kube, cfg.NotifLabel, sts.Namespace, vol.PersistentVolumeClaim.ClaimName)
		}

	}
}

// triggered on sts deletion event
// will add labels to associated pvcs

func handleDeleted(kube kubernetes.Interface, cfg structInternal.ControllerConfig, sts *appsv1.StatefulSet) {
	log.Printf("[INFO] STS deleted: %s", sts.Name)

	for _, vol := range sts.Spec.Template.Spec.Volumes {
		if vol.PersistentVolumeClaim == nil {
			continue
		}

		// get pvc object to check storage class
		pvcObj, err := kube.CoreV1().PersistentVolumeClaims(sts.Namespace).Get(context.TODO(), vol.PersistentVolumeClaim.ClaimName, metav1.GetOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to find PVC object %s: %s", vol.PersistentVolumeClaim.ClaimName, err)
			continue
		}

		log.Printf("[INFO] Found PVC object %s", pvcObj.Name)

		// ignore if storage class not in config
		if IgnoreStorageClass(pvcObj.Spec.StorageClassName, cfg.StorageClass) {
			continue
		}

		log.Printf("[INFO] Adding labels.")
		SetPvcLabel(kube, cfg.TimeLabel, time.Now().Format(cfg.TimeFormat), sts.Namespace, vol.PersistentVolumeClaim.ClaimName)
		SetPvcLabel(kube, cfg.NotifLabel, "0", sts.Namespace, vol.PersistentVolumeClaim.ClaimName)
	}
}

func IgnoreStorageClass(name *string, storageClasses []string) bool {
	if len(storageClasses) == 0 {
		return false
	}
	if name == nil {
		return !slices.Contains(storageClasses, "")
	}
	return !slices.Contains(storageClasses, *name)
}
