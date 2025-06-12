package kubernetes

import (
	"log"
	structInternal "volume-cleaner/internal/structure"

	"k8s.io/client-go/kubernetes"
)

func FindStale(kube kubernetes.Interface, cfg structInternal.SchedulerConfig) {
	log.Print("Job done")
}
