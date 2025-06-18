package main

import (
	// Standard Packages
	"log"
	"os"
	"strconv"

	// External Packages
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// Internal Packages
	kubeInternal "volume-cleaner/internal/kubernetes"
	structInternal "volume-cleaner/internal/structure"
)

func main() {
	log.Print("Volume cleaner scheduler started.")

	cfg := structInternal.SchedulerConfig{
		Namespace:   os.Getenv("NAMESPACE"),
		Label:       os.Getenv("LABEL"),
		TimeFormat:  os.Getenv("TIME_FORMAT"),
		GracePeriod: parseGracePeriod(os.Getenv("GRACE_PERIOD")),
		DryRun:      os.Getenv("DRY_RUN") == "true" || os.Getenv("DRY_RUN") == "1",
	}

	kubeClient, err := initKubeClient()
	if err != nil {
		log.Fatalf("Error creating kube client: %s", err)
	}

	kubeInternal.FindStale(kubeClient, cfg)
}

// read grace period value provided in the config and convert it to an int

func parseGracePeriod(value string) int {
	// Atoi means ASCII to Integer

	days, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Error parsing grace period value: %s", err)
	} else if days < 1 {
		log.Fatal("For safety reasons, grace period cannot be lower than one day.")
	}
	return days
}

// go client used to interact with k8s clusters

func initKubeClient() (*kubernetes.Clientset, error) {
	// service runs inside cluster as a pod, therefore will use in-cluster config
	// to connect with kubernetes API

	cfg, err := rest.InClusterConfig()
	if err == nil {
		return kubernetes.NewForConfig(cfg)
	}
	return nil, err
}
