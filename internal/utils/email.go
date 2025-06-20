package utils

import (
	// standard packages
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	// internal packages
	structInternal "volume-cleaner/internal/structure"
)

func SendNotif(client *http.Client, conf structInternal.EmailConfig, email string, personal structInternal.Personalisation) bool {

	url := conf.BaseURL + conf.Endpoint

	// Request Body
	reqBody, err := json.Marshal(
		structInternal.RequestBody{
			EmailAddress:    email,
			TemplateID:      conf.EmailTemplateID,
			Personalisation: personal,
		})

	if err != nil {
		log.Fatalf("Error creating request body: %v", err)
	}

	// Create the request and add the required headers
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Add("Authorization", "ApiKey-v1 "+conf.APIKey)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Send Request
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error making HTTP POST request: %v", err)

		// sending the email failed, but don't stop the program
		return false
	}
	log.Printf("Successfully Sent Email Notif to %s: %s", personal.Name, resp.Status)

	defer resp.Body.Close()

	return resp.StatusCode != 201 // return err boolean
}

// NOTE: This is what you would do before you call the SendNotif function (DELETE THIS WHEN THIS FUNCTION IS INTEGRATED WITH THE SCHEDULER)
// func example() {
// 	// setup
// 	email := "simulate-delivered@notification.canada.ca"
//
// 	personal := structInternal.Personalisation{
// 		Name:         "John Doe",
// 		VolumeName:   "Volume",
// 		VolumeID:     "Volume ID",
// 		GracePeriod:  "180", // in days
// 		DeletionDate: "June 17, 2025",
// 	}
//
// 	config := structInternal.EmailConfig{
// 		BaseURL:         "https://api.notification.canada.ca",
// 		Endpoint:        "/v2/notifications/email",
// 		EmailTemplateID: "Random Template",
// 		APIKey:          "Random APIKEY",
// 	}
//
// 	client := &http.Client{Timeout: 10 * time.Second}
//
// 	// sending email!
// 	code := SendNotif(client, config, email, personal)
//
// 	log.Printf("Status: %t", code)
// }
