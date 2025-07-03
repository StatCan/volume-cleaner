package kubernetes

import (
	// standard packages
	"context"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	/* Unfortunate that a lot of the kubernetes packages require renaming because
	they do not abide by good package name conventions as per https://go.dev/blog/package-names
	*/

	// external packages
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	// internal packages
	structInternal "volume-cleaner/internal/structure"
	utilsInternal "volume-cleaner/internal/utils"
)

// will write unit test for this when whole function is done

func FindStale(kube kubernetes.Interface, cfg structInternal.SchedulerConfig) {
	// One http client is created for emailing users
	client := &http.Client{Timeout: 10 * time.Second}
	errCount := 0

	deleteCount := 0
	emailCount := 0

	// Sort in descending order
	sort.Slice(cfg.NotifTimes, func(i, j int) bool {
		return cfg.NotifTimes[i] > cfg.NotifTimes[j] // Descending
	})

	for _, pvc := range PvcList(kube, cfg.Namespace) {
		log.Printf("Found pvc %s from namespace %s", pvc.Name, pvc.Namespace)

		// check if label exists (meaning unattached)
		// if pvc is attached to a sts, it would've had its label removed by the controller

		/*
			Even though having many nested if strctures goes against go style,
			because this code is so consequential (deleting volumes),
			I wanted to keep the logic here as straightfoward as possible
		*/

		timestamp, ok := pvc.Labels[cfg.Label]
		if ok {
			stale, staleError := IsStale(timestamp, cfg.TimeFormat, cfg.GracePeriod)
			if staleError != nil {
				log.Printf("Could not parse time: %s", staleError)
				errCount++
				continue
			}

			if stale {
				if cfg.DryRun {

					log.Printf("DRY RUN: delete pvc %s", pvc.Name)
					continue

				}

				err := kube.CoreV1().PersistentVolumeClaims(pvc.Namespace).Delete(context.TODO(), pvc.Name, metav1.DeleteOptions{})
				if err != nil {
					log.Printf("Error deleting pvc %s: %s", pvc.Name, err)
					errCount++
				}
				log.Print("PVC successfully deleted.")
				deleteCount++

			} else {
				log.Print("Grace period not passed.")

				notifCount, ok := pvc.Labels["volume-cleaner/notification-count"]
				if !ok {
					// TODO: replace this with variable, abstracted out label
					log.Print("Error reading label: volume-cleaner/notification-count")
					continue
				}
				currNotif, err := strconv.Atoi(notifCount)

				if err != nil {
					log.Printf("Error converting string label to int: %s", notifCount)
				}

				shouldSend, mailError := ShouldSendMail(timestamp, currNotif, cfg)

				if mailError != nil {
					log.Printf("Could not parse time: %s", mailError)
					errCount++
					continue
				}

				if shouldSend {
					if cfg.DryRun {
						log.Print("DRY RUN: email user")
						continue
					}

					// personal consists of details passed into the email template as variables while email is the email address that is consistent regardless of the template
					email, personal := utilsInternal.EmailDetails(kube, pvc, cfg.GracePeriod)

					err := utilsInternal.SendNotif(client, cfg.EmailCfg, email, personal)

					if err != nil {
						log.Printf("Error: Unable to send an email to %s at %s", personal.Name, email)
						errCount++
					} else {
						// Update Email Count
						emailCount++

						// Increment notification count by 1
						newNotifCount := strconv.Itoa(currNotif + 1)
						SetPvcLabel(kube, "volume-cleaner/notification-count", newNotifCount, pvc.Namespace, pvc.Name)
					}

				}
			}
		} else {
			log.Print("Not labelled. Skipping.")
		}
	}

	log.Printf("Job errors %d", errCount)
	log.Printf("Emails sent: %d", emailCount)
	log.Printf("Pvcs deleted: %d", deleteCount)
}

// determines if the grace period is greater than a given timestamp

func IsStale(timestamp string, format string, gracePeriod int) (bool, error) {
	timeObj, err := time.Parse(format, timestamp)
	if err != nil {
		return false, err
	}

	// difference in days
	diff := time.Since(timeObj).Hours() / 24

	log.Printf("Parsed timestamp: %f days.", diff)

	stale := int(diff) > gracePeriod

	log.Printf("int(diff) > cfg.GracePeriod: %v > %v == %v", int(diff), gracePeriod, stale)

	return stale, nil
}

func ShouldSendMail(timestamp string, currNotif int, cfg structInternal.SchedulerConfig) (bool, error) {
	log.Print("Checking email times....")

	timeObj, err := time.Parse(cfg.TimeFormat, timestamp)
	if err != nil {
		return false, err
	}
	daysLeft := cfg.GracePeriod - int(math.Floor(time.Since(timeObj).Hours()/24))

	log.Printf("Days left until deletion: %d", daysLeft)

	if currNotif < len(cfg.NotifTimes) && cfg.NotifTimes[currNotif] >= daysLeft {
		return true, nil
	}

	return false, nil
}
