package main

import (
	// standard Packages
	"context"
	"log"
	"os"

	// external Packages

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// internal Packages
	kubeInternal "volume-cleaner/internal/kubernetes"
	structInternal "volume-cleaner/internal/structure"
)

func main() {
	/*
		It took me a while to figure this out because there wasn't much documentation about this

		log.Println and log.Print can be *practically* used for the purpose

		Both functions actually make calls to log.Output. The only difference lies in how the string
		is formatted before printing. log.Println uses fmt.Appendln and log.Print uses fmt.Append. The
		reason why they can be used interchangebly is because log.Ouput automatically inserts a
		new line character if it's not present
	*/

	log.Print("Volume cleaner controller started.")

	cfg := structInternal.ControllerConfig{
		Namespace:  os.Getenv("NAMESPACE"),
		Label:      os.Getenv("LABEL"),
		TimeFormat: os.Getenv("TIME_FORMAT"),
	}

	kubeClient, err := initKubeClient()
	if err != nil {
		// log.Fatalf will automatically call os.Exit

		log.Fatalf("Error creating kube client: %s", err)
	}

	// scans pvcs to find already unattached ones
	kubeInternal.InitialScan(kubeClient, cfg)

	// watches stateful sets to discover newly unattached pvcs
	kubeInternal.WatchSts(context.TODO(), kubeClient, cfg)
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
