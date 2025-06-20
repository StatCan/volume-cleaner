package utils

import (
	// standard packages
	"log"
	"net/http"
	"testing"
	"time"

	// external packages
	"github.com/stretchr/testify/assert"

	// internal packages
	structInternal "volume-cleaner/internal/structure"
)

func TestSendingNotif(t *testing.T) {
	email := "simulate-delivered@notification.canada.ca"

	personal := structInternal.Personalisation{
		Name:         "John Doe",
		VolumeName:   "Volume",
		VolumeID:     "Volume ID",
		GracePeriod:  "180", // in days
		DeletionDate: "June 17, 2025",
	}

	client := &http.Client{Timeout: 10 * time.Second}

	configInvalid := structInternal.EmailConfig{
		BaseURL:         "https://api.notification.canada.ca",
		Endpoint:        "/v2/notifications/email",
		EmailTemplateID: "Random Template",
		APIKey:          "Random Key",
	}

	// sending email!
	err := SendNotif(client, configInvalid, email, personal)

	log.Printf("Status: %t", err)

	t.Run("sending an unauthorized api email request", func(t *testing.T) {
		assert.Equal(t, err, true)
	})
}
