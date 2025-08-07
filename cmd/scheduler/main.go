package main

import (
	// standard Packages
	"log"
	"os"

	// internal Packages
	kubeInternal "volume-cleaner/internal/kubernetes"
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

	// Scheduler struct which composes an EmailConfig
	// there is also a config for the controller
	cfg := structInternal.SchedulerConfig{
		Namespace:   os.Getenv("NAMESPACE"),
		TimeLabel:   os.Getenv("TIME_LABEL"),
		NotifLabel:  os.Getenv("NOTIF_LABEL"),
		TimeFormat:  os.Getenv("TIME_FORMAT"),
		GracePeriod: utilsInternal.ParseGracePeriod(os.Getenv("GRACE_PERIOD")),
		DryRun:      os.Getenv("DRY_RUN") == "true" || os.Getenv("DRY_RUN") == "1",
		NotifTimes:  utilsInternal.ParseNotifTimes(os.Getenv("NOTIF_TIMES")),
		EmailCfg:    emailCfg,
	}

	// init client to interact with k8s cluster
	kubeClient, err := kubeInternal.InitKubeClient()
	if err != nil {
		log.Fatalf("Error creating kube client: %s", err)
	}

	// run main scheduler logic
	kubeInternal.FindStale(kubeClient, cfg)
}
