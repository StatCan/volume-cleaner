package utils

import (
	// standard packages
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	// external packages
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	// internal packages
	structInternal "volume-cleaner/internal/structure"
)

// given a collection of configs, this function makes a post request to an third party email service and sends an email to a user

func SendNotif(client *http.Client, conf structInternal.EmailConfig, email string, personal structInternal.Personalisation) error {

	url := conf.BaseURL + conf.Endpoint

	// Request Body
	reqBody, err := json.Marshal(
		structInternal.RequestBody{
			EmailAddress:    email,
			TemplateID:      conf.EmailTemplateID,
			Personalisation: personal,
		})

	if err != nil {
		log.Printf("[ERROR] Failed to create request body: %v", err)
	}

	// Create the request and add the required headers
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	request.Header.Add("Authorization", "ApiKey-v1 "+conf.APIKey)
	request.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Printf("[ERROR] Failed to create request: %v", err)
	} else {
		log.Printf("[INFO] Successfully created HTTP request.")
	}

	// Send Request
	response, err := client.Do(request)

	if response != nil {
		defer response.Body.Close()
	}

	log.Print(personal.DeletionDate)
	log.Print(personal.GracePeriod)
	log.Print(personal.Name)
	log.Print(personal.VolumeName)

	if request != nil {
		log.Print(request.Header)
	}

	if err != nil {
		log.Printf("[ERROR] Failed to create HTTP POST request: %v", err)

		// sending the email failed, but don't stop the program
		return errors.New("error response is invalid")
	}

	if response.StatusCode == 201 {
		log.Printf("[INFO] Successfully sent email notification to %s: %s", personal.Name, response.Status)

		return nil
	}
	return errors.New(response.Status)
}

// given a pvc, this function will aquire the details related to the pvc such as the owner of the pvc, their email, the bounded volume name and ID, and details about its deletion

func EmailDetails(kube kubernetes.Interface, pvc corev1.PersistentVolumeClaim, gracePeriod int) (string, structInternal.Personalisation) {
	ns := pvc.Namespace
	vol := pvc.Spec.VolumeName

	// Acquire User Email
	email := nsEmail(kube, ns)

	// Calculate DeletionDate
	now := time.Now()
	futureTime := now.Add(time.Duration(gracePeriod) * 24 * time.Hour)

	personal := structInternal.Personalisation{
		Name:         ns,
		VolumeName:   vol,
		GracePeriod:  fmt.Sprintf("%d", gracePeriod),
		DeletionDate: futureTime.Format(time.UnixDate),
	}

	return email, personal
}

// returns the email associated with a namespace

func nsEmail(kube kubernetes.Interface, name string) string {
	ns, err := kube.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] Failed to get namespace %s: %v", name, err)
	}

	email := ns.Annotations["owner"]
	if email == "" {
		log.Printf("[ERROR] Annotation 'owner' for namespace %s is empty", name)
	} else {
		log.Printf("[INFO] Successfully acquired owner email %s", email)
	}

	return email
}
