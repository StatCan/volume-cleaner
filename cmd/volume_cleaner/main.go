package main

import (
	// External Packages
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// Internal Packages
	kubeInternal "volume-cleaner/internal/kubernetes"
	kubeStructure "volume-cleaner/internal/structure"
)

func main() {

	fmt.Println("Volume cleaner started.")

	cfg := kubeStructure.Config{
		Namespace:  os.Getenv("NAMESPACE"),
		Label:      os.Getenv("LABEL"),
		TimeFormat: os.Getenv("TIME_FORMAT"),
	}

	kubeClient, err := initKubeClient()
	if err != nil {
		log.Fatalf("Error creating kube client: %v", err)
	}

	for _, pvc := range kubeInternal.FindUnattachedPVCs(kubeClient) {
		kubeInternal.SetPvcLabel(kubeClient, cfg.Label, time.Now().Format(cfg.TimeFormat), pvc.Namespace, pvc.Name)
	}

	kubeInternal.WatchSts(context.TODO(), kubeClient, cfg)

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
