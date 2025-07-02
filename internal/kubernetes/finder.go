package kubernetes

import (
	// standard packages
	"context"
	"log"
	"math"
	"net/http"
	"time"

	/* Unfortunate that a lot of the kubernetes packages require renaming because
	they do not abide by good package name conventions as per https://go.dev/blog/package-names
	*/

	// external packages
	corev1 "k8s.io/api/core/v1"
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

			} else {
				log.Print("Grace period not passed.")

				shouldSend, mailError := ShouldSendMail(timestamp, pvc, cfg)

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
					}

				}
			}
		} else {
			log.Print("Not labelled. Skipping.")
		}
	}

	log.Printf("Job errors %d", errCount)
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

func ShouldSendMail(timestamp string, _ corev1.PersistentVolumeClaim, cfg structInternal.SchedulerConfig) (bool, error) {
	log.Print("Checking email times....")

	timeObj, err := time.Parse(cfg.TimeFormat, timestamp)
	if err != nil {
		return false, err
	}
	daysLeft := cfg.GracePeriod - int(math.Floor(time.Since(timeObj).Hours()/24))

	log.Printf("Days left until deletion: %d", daysLeft)

	for _, time := range cfg.NotifTimes {
		if daysLeft == time {
			return true, nil
		}
	}

	return false, nil
}
