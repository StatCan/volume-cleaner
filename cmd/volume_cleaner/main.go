package main

import (
	// External Packages
	"fmt"
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// Internal Packages
	kuber "volume-cleaner/internal/kubernetes"
)

type Config struct {
	DryRun      bool
	GracePeriod int
}

func main() {
	fmt.Println("Volume cleaner started.")

	cfg := Config{
		DryRun:      os.Getenv("DRY_RUN") == "true",
		GracePeriod: 30,
	}

	kubeClient, err := initKubeClient()
	if err != nil {
		log.Fatalf("Error creating kube client: %v", err)
	}

	cleanVolumes(kubeClient, cfg)
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

func cleanVolumes(kube kubernetes.Interface, cfg Config) {
	kuber.FindUnattachedPVCs(kube)
}
