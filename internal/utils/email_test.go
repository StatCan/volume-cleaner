package utils

import (
	// standard packages
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	// external packages
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	// internal packages
	structInternal "volume-cleaner/internal/structure"
)

func TestSendingNotif(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	email := "simulate-delivered@notification.canada.ca"

	personal := structInternal.Personalisation{
		Name:         "John Doe",
		VolumeName:   "Volume",
		VolumeID:     "Volume ID",
		GracePeriod:  "180", // in days
		DeletionDate: "June 17, 2025",
	}

	configValid := structInternal.EmailConfig{
		BaseURL:         "https://api.notification.canada.ca",
		Endpoint:        "/v2/notifications/email",
		EmailTemplateID: os.Getenv("EMAIL_TEMPLATE_ID"),
		APIKey:          os.Getenv("GC_NOTIFY_API_KEY_TEST"),
	}

	val := os.Getenv("EMAIL_TEMPLATE_ID")

	log.Printf("id: %s", val)

	client := &http.Client{Timeout: 10 * time.Second}

	// sending email!
	success := sendNotif(client, configValid, email, personal)

	log.Printf("Status: %t", success)

	t.Run("sending an authorized api email request", func(t *testing.T) {
		assert.Equal(t, success, true)
	})

	configInvalid := structInternal.EmailConfig{
		BaseURL:         "https://api.notification.canada.ca",
		Endpoint:        "/v2/notifications/email",
		EmailTemplateID: "Random Template",
		APIKey:          "Random Key",
	}

	// sending email!
	fail := sendNotif(client, configInvalid, email, personal)

	log.Printf("Status: %t", fail)

	t.Run("sending an unauthorized api email request", func(t *testing.T) {
		assert.Equal(t, fail, false)
	})
}
