package main

import (
	// External Packages

	"log"
	"os"
	"strconv"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// Internal Packages
	kubeInternal "volume-cleaner/internal/kubernetes"
	structInternal "volume-cleaner/internal/structure"
)

func main() {

	log.Print("Volume cleaner scheduler started!")

	cfg := structInternal.SchedulerConfig{
		Namespace:   os.Getenv("NAMESPACE"),
		Label:       os.Getenv("LABEL"),
		TimeFormat:  os.Getenv("TIME_FORMAT"),
		GracePeriod: parseGracePeriod(os.Getenv("GRACE_PERIOD")),
		DryRun:      os.Getenv("DRY_RUN") == "true" || os.Getenv("DRY_RUN") == "1",
	}

	kubeClient, err := initKubeClient()
	if err != nil {
		log.Fatalf("Error creating kube client: %v", err)
	}

	kubeInternal.FindStale(kubeClient, cfg)

}

func parseGracePeriod(value string) int {
	days, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Error parsing grace period value: %v", err)
	} else if days < 1 {
		log.Fatal("For saftey reasons, grace period cannot be lower than one day.")
	}
	return days
}

func initKubeClient() (*kubernetes.Clientset, error) {
	// service runs inside cluster as a pod, therefore will use in-cluster config
	// to connect with kubernetes API

	cfg, err := rest.InClusterConfig()
	if err == nil {
		return kubernetes.NewForConfig(cfg)
	}
	return nil, err
}
