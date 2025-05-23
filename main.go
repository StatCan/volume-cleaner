package main

import (
	"context"
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

func main() {
	fmt.Println("Volume cleaner started.")

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

// pointers?

func initKubeClient() (*kubernetes.Clientset, error) {
	// uses in-cluster config
	cfg, err := rest.InClusterConfig()
	if err == nil {
		return kubernetes.NewForConfig(cfg)
	}
	return nil, err
}

func initGraphClient() (*msgraphsdk.GraphServiceClient, error) {
	// currently just mocks an empty client
	return &msgraphsdk.GraphServiceClient{}, nil
}

func cleanVolumes(ctx context.Context, kube kubernetes.Interface, graph *msgraphsdk.GraphServiceClient, cfg Config) {
	ns, err := kube.CoreV1().Namespaces().List(ctx, metav1.ListOptions{LabelSelector: "kubernetes.io/metadata.name=das"})
	if err != nil {
		log.Fatalf("Error listing volumes: %v", err)
	}

	fmt.Println(ns)

	vols, err := kube.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing volumes: %v", err)
	}

	fmt.Println(vols)

}
