package main

import (
	// External Packages
	"context"
	"fmt"
	"log"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// Internal Packages
	kubeInternal "volume-cleaner/internal/kubernetes"
)

func main() {
	fmt.Println("Volume cleaner started.")

	kubeClient, err := initKubeClient()
	if err != nil {
		log.Fatalf("Error creating kube client: %v", err)
	}

	for _, pvc := range kubeInternal.FindUnattachedPVCs(kubeClient) {
		kubeInternal.SetPvcLabel(kubeClient, "volume-cleaner/unattached-time", time.Now().Format("2006-01-02_15-04-05Z"), pvc.Namespace, pvc.Name)
	}

	kubeInternal.WatchSts(context.TODO(), kubeClient, "anray-liu")

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
