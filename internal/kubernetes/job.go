package kubernetes

import (
	"log"
	structInternal "volume-cleaner/internal/structure"

	"k8s.io/client-go/kubernetes"
)

func FindStale(_ kubernetes.Interface, _ structInternal.SchedulerConfig) {
	log.Print("Job done")
}
