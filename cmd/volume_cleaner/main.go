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
	ClientID       string
	ClientSecret   string
	TenantID       string
	DryRun         bool
	AllowedDomains []string
	GracePeriod    int
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
	vols, err := kube.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing volumes: %v", err)
	}

	for _, vol := range vols.Items {
		fmt.Println(vol.Name, vol.Spec.ClaimRef.Name, vol.Status.LastPhaseTransitionTime.Time.String())
	}

}
