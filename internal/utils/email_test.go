package utils

import (
	// standard packages
	"errors"
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

// TestSendingNotif verifies that SendNotif returns the expected error when using invalid API credentials.
func TestSendingNotif(t *testing.T) {
	email := "simulate-delivered@notification.canada.ca"

	personal := structInternal.Personalisation{
		Name:         "John Doe",
		VolumeName:   "Volume",
		DaysLeft:     "180", // in days
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
		assert.Equal(t, err, errors.New("403 Forbidden"))
	})
}

// TestEmailDetails covers multiple scenarios for retrieving email and personalisation data from Kubernetes objects.
func TestEmailDetails(t *testing.T) {
	tests := []struct {
		name                    string
		namespace               *corev1.Namespace
		pvc                     corev1.PersistentVolumeClaim
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
			expectedEmail: "test@example.com",
			expectedPersonalisation: structInternal.Personalisation{
				Name:       "test-namespace",
				VolumeName: "test-pvc",
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
			expectedEmail: "", // Should be empty if owner annotation is missing
			expectedPersonalisation: structInternal.Personalisation{
				Name:       "no-owner-ns",
				VolumeName: "test-pvc-no-owner",
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
			expectedEmail: "", // Should be empty if namespace is not found
			expectedPersonalisation: structInternal.Personalisation{
				Name:       "non-existent-ns", // The name from PVC is used even if namespace isn't found
				VolumeName: "test-pvc-non-existent-ns",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.namespace != nil {
				kubeClient := fake.NewClientset(tt.namespace)

				email, personal := EmailDetails(kubeClient, tt.pvc, 0.0)

				// Assert the email
				assert.Equal(t, tt.expectedEmail, email, "Email should match")

				// Assert the Personalisation struct fields, handling the time dynamically
				assert.Equal(t, tt.expectedPersonalisation.Name, personal.Name, "Personalisation Name should match")
				assert.Equal(t, tt.expectedPersonalisation.VolumeName, personal.VolumeName, "Personalisation VolumeName should match")

			} else {
				// For the "Non-existent Namespace" case, create a client without the namespace
				kubeClient := fake.NewClientset()
				email, personal := EmailDetails(kubeClient, tt.pvc, 0.0)

				assert.Equal(t, tt.expectedEmail, email, "Email should be empty for non-existent namespace")
				assert.Equal(t, tt.expectedPersonalisation.Name, personal.Name, "Personalisation Name should match for non-existent namespace")
				assert.Equal(t, tt.expectedPersonalisation.VolumeName, personal.VolumeName, "Personalisation VolumeName should match for non-existent namespace")
			}
		})
	}
}
