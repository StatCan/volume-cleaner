package main

import (
	// standard Packages
	"context"
	"log"
	"os"

	// internal Packages
	kubeInternal "volume-cleaner/internal/kubernetes"
	structInternal "volume-cleaner/internal/structure"
)

func main() {
	/*
		It took me a while to figure this out because there wasn't much documentation about this

		log.Println and log.Print can be *practically* used for the same purpose

		Both functions actually make calls to log.Output. The only difference lies in how the string
		is formatted before printing. log.Println uses fmt.Appendln and log.Print uses fmt.Append. The
		reason why they can be used interchangebly is because log.Ouput automatically inserts a
		new line character if it's not present
	*/

	log.Print("[INFO] Volume cleaner controller started.")

	// controller config
	// there is also a config for the scheduler

	cfg := structInternal.ControllerConfig{
		Namespace:    os.Getenv("NAMESPACE"),
		TimeLabel:    os.Getenv("TIME_LABEL"),
		NotifLabel:   os.Getenv("NOTIF_LABEL"),
		TimeFormat:   os.Getenv("TIME_FORMAT"),
		StorageClass: os.Getenv("STORAGE_CLASS"),
	}

	// init client to interact with k8s cluster

	kubeClient, err := kubeInternal.InitKubeClient()
	if err != nil {
		// log.Fatalf will automatically call os.Exit

		log.Fatalf("[ERROR] Failed to create kube client: %s", err)
	}

	// scans pvcs to find already unattached pvcs
	kubeInternal.InitialScan(kubeClient, cfg)

	// watches stateful sets to discover newly unattached pvcs
	kubeInternal.WatchSts(context.TODO(), kubeClient, cfg)
}
