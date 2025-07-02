package main

import (
	// Standard Packages
	"log"
	"os"

	// External Packages
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// Internal Packages
	structInternal "volume-cleaner/internal/structure"
	utilsInternal "volume-cleaner/internal/utils"
)

func main() {
	log.Print("Volume cleaner scheduler started.")

	// Initialize an EmailConfig struct
	emailCfg := structInternal.EmailConfig{
		BaseURL:         os.Getenv("BASE_URL"),
		Endpoint:        os.Getenv("ENDPOINT"),
		EmailTemplateID: os.Getenv("EMAIL_TEMPLATE_ID"),
		APIKey:          os.Getenv("API_KEY"),
	}

	gracePeriod := utilsInternal.ParseGracePeriod(os.Getenv("GRACE_PERIOD"))
	notifTimes := utilsInternal.ParseNotifTimes(os.Getenv("NOTIF_TIMES"))

	cfg := structInternal.SchedulerConfig{
		Namespace:   os.Getenv("NAMESPACE"),
		Label:       os.Getenv("LABEL"),
		TimeFormat:  os.Getenv("TIME_FORMAT"),
		GracePeriod: gracePeriod,
		DryRun:      os.Getenv("DRY_RUN") == "true" || os.Getenv("DRY_RUN") == "1",
		NotifTimes:  notifTimes,
		EmailCfg:    emailCfg,
	}

	kubeClient, err := initKubeClient()
	if err != nil {
		log.Fatalf("Error creating kube client: %s", err)
	}

	FindStale(kubeClient, cfg)
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
