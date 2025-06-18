package kubernetes

import (
	"context"
	"log"
	"time"

	structInternal "volume-cleaner/internal/structure"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

// will write unit test for this when whole function is done

func FindStale(kube kubernetes.Interface, cfg structInternal.SchedulerConfig) {
	for _, pvc := range PvcList(kube, cfg.Namespace) {
		log.Printf("Found a pvc: %s from namespace %s", pvc.Name, pvc.Namespace)
		timestamp, ok := pvc.Labels[cfg.Label]
		if ok {
			if IsStale(timestamp, cfg.TimeFormat, cfg.GracePeriod) {
				if cfg.DryRun {
					log.Printf("DRY RUN: delete pvc %s", pvc.Name)
				} else {
					err := kube.CoreV1().PersistentVolumeClaims(pvc.Namespace).Delete(context.TODO(), pvc.Name, metav1.DeleteOptions{})
					if err != nil {
						log.Fatalf("Error deleting pvc %v: %v", pvc.Name, err)
					}
					log.Print("PVC successfully deleted.")
				}
			} else {
				log.Print("Grace period not passed. Skipping.")
			}
		} else {
			log.Print("Not labelled. Skipping.")
		}
	}
}

func IsStale(timestamp string, format string, gracePeriod int) bool {
	timeObj, err := time.Parse(format, timestamp)
	if err != nil {
		log.Fatalf("Could not parse time: %s", err)
	}

	diff := time.Since(timeObj).Hours() / 24

	log.Printf("Parsed timestamp: %f days.", diff)

	stale := int(diff) > gracePeriod

	log.Printf("int(diff) > cfg.GracePeriod: %v > %v == %v", int(diff), gracePeriod, stale)

	return stale
}
