package main

import (
	"context"
	"fmt"
	"log"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// will eventually move to config
const GRACE_PERIOD = 30

// pointers?
func initKubeClient() (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(&rest.Config{
		Host: "",
	})
}

func initGraphClient() (*msgraphsdk.GraphServiceClient, error) {
	return &msgraphsdk.GraphServiceClient{}, nil
}

func cleanVolumes(ctx context.Context, kube kubernetes.Interface, graph *msgraphsdk.GraphServiceClient) {
	fmt.Println(graph)
	fmt.Println(kube)

	volumesList, err := kube.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing volumes: %v", err)
	}

	fmt.Println(volumesList)

}

func main() {
	kubeClient, err := initKubeClient()
	if err != nil {
		log.Fatalf("Error creating kube client: %v", err)
	}

	graphClient, err := initGraphClient()
	if err != nil {
		log.Fatalf("Error creating graph client: %v", err)
	}

	cleanVolumes(context.Background(), kubeClient, graphClient)

}
