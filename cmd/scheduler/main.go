package main

import (
	// Standard Packages
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	// External Packages
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// Internal Packages
	kubeInternal "volume-cleaner/internal/kubernetes"
	structInternal "volume-cleaner/internal/structure"
)

func main() {
	log.Print("Volume cleaner scheduler started.")

	cfg := structInternal.SchedulerConfig{
		Namespace:   os.Getenv("NAMESPACE"),
		Label:       os.Getenv("LABEL"),
		TimeFormat:  os.Getenv("TIME_FORMAT"),
		GracePeriod: ParseGracePeriod(os.Getenv("GRACE_PERIOD")),
		DryRun:      os.Getenv("DRY_RUN") == "true" || os.Getenv("DRY_RUN") == "1",
		NotifTimes:  ParseNotifTimes(os.Getenv("NOTIF_TIMES")),
	}

	kubeClient, err := initKubeClient()
	if err != nil {
		log.Fatalf("Error creating kube client: %s", err)
	}

	kubeInternal.FindStale(kubeClient, cfg)
}

func ParseNotifTimes(str string) []int {
	var intSlice []int

	// use fields() and join() to get rid of all whitespace
	// split by delimeter ,
	// try to convert each value to an int, error out if failed
	// sort final slice of ints

	parsedString := strings.Split(strings.Join(strings.Fields(str), ""), ",")
	for _, val := range parsedString {
		converted, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("Error parsing notification time: %s", err)
		}
		intSlice = append(intSlice, converted)
	}

	sort.Ints(intSlice)

	log.Print(intSlice)

	return intSlice
}

// read grace period value provided in the config and convert it to an int

func ParseGracePeriod(value string) int {
	// Atoi means ASCII to Integer

	days, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Error parsing grace period value: %s", err)
	} else if days < 1 {
		log.Fatal("For safety reasons, grace period cannot be lower than one day.")
	}
	return days
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
