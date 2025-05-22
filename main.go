package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
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

// pointers?

func initKubeClient() (*kubernetes.Clientset, error) {

	// Try in-cluster config (only available inside Kubernetes pods)
	if inClusterCfg, err := rest.InClusterConfig(); err == nil {
		return kubernetes.NewForConfig(inClusterCfg)
	}

	// Fall back to mockable out-of-cluster config if KUBECONFIG is set
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubernetes.NewForConfig(&rest.Config{
			Host: "http://localhost:8080",
		})
	}

	// Default failure case
	return nil, errors.New("no valid Kubernetes config found")
}

func initGraphClient() (*msgraphsdk.GraphServiceClient, error) {
	return &msgraphsdk.GraphServiceClient{}, nil
}

func cleanVolumes(ctx context.Context, kube kubernetes.Interface, graph *msgraphsdk.GraphServiceClient, cfg Config) {
	fmt.Println(graph)
	fmt.Println(kube)

	volumesList, err := kube.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing volumes: %v", err)
	}

	fmt.Println(volumesList)

}

func main() {
	cfg := Config{
		ClientID:       os.Getenv("CLIENT_ID"),
		ClientSecret:   os.Getenv("CLIENT_SECRET"),
		TenantID:       os.Getenv("TENANT_ID"),
		DryRun:         os.Getenv("DRY_RUN") == "true",
		AllowedDomains: strings.Split(os.Getenv("ALLOWED_DOMAINS"), ","),
		GracePeriod:    30,
	}

	kubeClient, err := initKubeClient()
	if err != nil {
		log.Fatalf("Error creating kube client: %v", err)
	}

	graphClient, err := initGraphClient()
	if err != nil {
		log.Fatalf("Error creating graph client: %v", err)
	}

	cleanVolumes(context.Background(), kubeClient, graphClient, cfg)

}
