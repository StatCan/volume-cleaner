package kubernetes

import (
	"log"
	"time"
	structInternal "volume-cleaner/internal/structure"

	"k8s.io/client-go/kubernetes"
)

func FindStale(kube kubernetes.Interface, cfg structInternal.SchedulerConfig) {
	for _, pvc := range PvcList(kube, cfg.Namespace) {
		log.Printf("Found a pvc: %s from namespace %s", pvc.Name, pvc.Namespace)
		timestamp, ok := pvc.Labels[cfg.Label]
		if ok {
			time_obj, err := time.Parse(cfg.TimeFormat, timestamp)
			if err != nil {
				log.Fatalf("Could not parse time: %s", err)
			}

			diff := time.Now().Sub(time_obj)

			log.Printf("This PVC is %f hours old.", diff.Hours())

			if cfg.DryRun {

			}

		} else {
			log.Print("Not labelled. Skipping.")
		}
	}
}
