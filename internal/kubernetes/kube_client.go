package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// go client used to interact with k8s clusters

func InitKubeClient() (*kubernetes.Clientset, error) {
	// service runs inside cluster as a pod, therefore will use in-cluster config
	// to connect with kubernetes API

	cfg, err := rest.InClusterConfig()
	if err == nil {
		return kubernetes.NewForConfig(cfg)
	}
	return nil, err
}
