package utils

import (
	// standard packages
	"log"
	"net/http"
	"testing"
	"time"

	// external packages
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	// internal packages
	structInternal "volume-cleaner/internal/structure"
)

func TestSendingNotif(t *testing.T) {
	email := "simulate-delivered@notification.canada.ca"

	personal := structInternal.Personalisation{
		Name:         "John Doe",
		VolumeName:   "Volume",
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

func TestEmailDetails(t *testing.T) {
	tests := []struct {
		name                    string
		namespace               *corev1.Namespace
		pvc                     corev1.PersistentVolumeClaim
		gracePeriod             int
		expectedEmail           string
		expectedPersonalisation structInternal.Personalisation
	}{
		{
			name: "Successful EmailDetails retrieval",
			namespace: &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-namespace",
					Annotations: map[string]string{
						"owner": "test@example.com",
					},
				},
			},
			pvc: corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pvc",
					Namespace: "test-namespace",
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					VolumeName: "pv-test-volume-123",
				},
			},
			gracePeriod:   7,
			expectedEmail: "test@example.com",
			expectedPersonalisation: structInternal.Personalisation{
				Name:        "test-namespace",
				VolumeName:  "pv-test-volume-123",
				GracePeriod: "7",
				// DeletionDate will be checked dynamically due to time.Now()
			},
		},
		{
			name: "Namespace without owner annotation",
			namespace: &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "no-owner-ns",
					Annotations: map[string]string{
						"some-other-annotation": "value",
					},
				},
			},
			pvc: corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pvc-no-owner",
					Namespace: "no-owner-ns",
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					VolumeName: "pv-no-owner-volume",
				},
			},
			gracePeriod:   14,
			expectedEmail: "", // Should be empty if owner annotation is missing
			expectedPersonalisation: structInternal.Personalisation{
				Name:        "no-owner-ns",
				VolumeName:  "pv-no-owner-volume",
				GracePeriod: "14",
				// DeletionDate will be checked dynamically
			},
		},
		{
			name:      "Non-existent Namespace",
			namespace: nil, // Simulate a non-existent namespace by not adding it to the fake client
			pvc: corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pvc-non-existent-ns",
					Namespace: "non-existent-ns",
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					VolumeName: "pv-non-existent-volume",
				},
			},
			gracePeriod:   3,
			expectedEmail: "", // Should be empty if namespace is not found
			expectedPersonalisation: structInternal.Personalisation{
				Name:        "non-existent-ns", // The name from PVC is used even if namespace isn't found
				VolumeName:  "pv-non-existent-volume",
				GracePeriod: "3",
				// DeletionDate will be checked dynamically
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.namespace != nil {
				kubeClient := fake.NewClientset(tt.namespace)

				email, personal := EmailDetails(kubeClient, tt.pvc, tt.gracePeriod)

				// Assert the email
				assert.Equal(t, tt.expectedEmail, email, "Email should match")

				// Assert the Personalisation struct fields, handling the time dynamically
				assert.Equal(t, tt.expectedPersonalisation.Name, personal.Name, "Personalisation Name should match")
				assert.Equal(t, tt.expectedPersonalisation.VolumeName, personal.VolumeName, "Personalisation VolumeName should match")
				assert.Equal(t, tt.expectedPersonalisation.GracePeriod, personal.GracePeriod, "Personalisation GracePeriod should match")

				// Calculate the expected deletion date based on the current time and grace period
				now := time.Now()
				futureTime := now.Add(time.Duration(tt.gracePeriod) * 24 * time.Hour)
				expectedDeletionDate := futureTime.Format(time.UnixDate)

				// Assert deletion date, allowing for a small time difference
				parsedActualDate, err := time.Parse(time.UnixDate, personal.DeletionDate)
				assert.NoError(t, err, "Should be able to parse actual DeletionDate")

				parsedExpectedDate, err := time.Parse(time.UnixDate, expectedDeletionDate)
				assert.NoError(t, err, "Should be able to parse expected DeletionDate")

				// Check if the difference is within an acceptable margin (e.g., 1 second)
				assert.WithinDuration(t, parsedExpectedDate, parsedActualDate, 1*time.Second, "DeletionDate should be approximately correct")

			} else {
				// For the "Non-existent Namespace" case, create a client without the namespace
				kubeClient := fake.NewClientset()
				email, personal := EmailDetails(kubeClient, tt.pvc, tt.gracePeriod)

				assert.Equal(t, tt.expectedEmail, email, "Email should be empty for non-existent namespace")
				assert.Equal(t, tt.expectedPersonalisation.Name, personal.Name, "Personalisation Name should match for non-existent namespace")
				assert.Equal(t, tt.expectedPersonalisation.VolumeName, personal.VolumeName, "Personalisation VolumeName should match for non-existent namespace")
				assert.Equal(t, tt.expectedPersonalisation.GracePeriod, personal.GracePeriod, "Personalisation GracePeriod should match for non-existent namespace")

				now := time.Now()
				futureTime := now.Add(time.Duration(tt.gracePeriod) * 24 * time.Hour)
				expectedDeletionDate := futureTime.Format(time.UnixDate)

				parsedActualDate, err := time.Parse(time.UnixDate, personal.DeletionDate)
				assert.NoError(t, err, "Should be able to parse actual DeletionDate")

				parsedExpectedDate, err := time.Parse(time.UnixDate, expectedDeletionDate)
				assert.NoError(t, err, "Should be able to parse expected DeletionDate")
				assert.WithinDuration(t, parsedExpectedDate, parsedActualDate, 1*time.Second, "DeletionDate should be approximately correct for non-existent namespace")
			}
		})
	}
}
