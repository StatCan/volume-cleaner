package main

import (
	"context"
	"fmt"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

// pointers?

func initKubeClient() (*kubernetes.Clientset, error) {
	// service runs inside cluster as a pod, therefore will use in-cluster config
	// to connect with cluster

	cfg, err := rest.InClusterConfig()
	if err == nil {
		return kubernetes.NewForConfig(cfg)
	}
	return nil, err
}

func cleanVolumes(kube kubernetes.Interface, cfg Config) {
	ns, err := kube.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing namespaces: %v", err)
	}

	for _, namespace := range ns.Items {
		fmt.Println(namespace.Name)

		pvcs, err := kube.CoreV1().PersistentVolumeClaims(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing volume claims: %v", err)
		}

		// azure disk will have the same name as the volume
		// e.g pvc-d88040d...

		for _, claim := range pvcs.Items {
			fmt.Println(claim.Name, claim.Spec.VolumeName)
		}

	}

	fmt.Println(ns.Items)

}
