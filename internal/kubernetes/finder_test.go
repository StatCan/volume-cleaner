package kubernetes

import (
	// standard packages
	"context"
	"sort"
	"testing"
	"time"

	// external packages
	"github.com/stretchr/testify/assert"

	// internal packages
	structInternal "volume-cleaner/internal/structure"
	testInternal "volume-cleaner/internal/tests"
)

func TestIsStale(t *testing.T) {

	t.Run("test successful determination of stale pvcs", func(t *testing.T) {
		format := "2006-01-02_15-04-05Z"

		type testCase struct {
			timestamp     string
			format        string
			gracePeriod   int
			expectedValue bool
		}

		testCases := []testCase{
			{
				timestamp:     time.Now().Add(-time.Hour * 24 * 181).Format(format),
				format:        format,
				gracePeriod:   180,
				expectedValue: true,
			},
			{
				timestamp:     time.Now().Add(-time.Hour*24*180 - time.Hour*23).Format(format),
				format:        format,
				gracePeriod:   180,
				expectedValue: false,
			},
			{
				timestamp:     time.Now().Add(-time.Hour*24*180 - time.Hour*23 - time.Minute*59 - time.Second*59).Format(format),
				format:        format,
				gracePeriod:   180,
				expectedValue: false,
			},
			{
				timestamp:     time.Now().Add(-time.Hour * 24 * 1000).Format(format),
				format:        format,
				gracePeriod:   180,
				expectedValue: true,
			},
			{
				timestamp:     time.Now().Format(format),
				format:        format,
				gracePeriod:   180,
				expectedValue: false,
			},
			{
				timestamp:     time.Now().Format(format),
				format:        format,
				gracePeriod:   0,
				expectedValue: false,
			},
			{
				timestamp:     time.Now().Add(-time.Second).Format(format),
				format:        format,
				gracePeriod:   0,
				expectedValue: false,
			},
			{
				timestamp:     time.Now().Add(-time.Hour * 24).Format(format),
				format:        format,
				gracePeriod:   0,
				expectedValue: true,
			},
		}

		for _, test := range testCases {
			v, err := IsStale(test.timestamp, test.format, test.gracePeriod)
			if err != nil {
				t.Fatal("IsStale failed.")
			}
			assert.Equal(t, v, test.expectedValue)
		}

	})
}

func TestShouldSendMail(t *testing.T) {

	t.Run("test successful determination of email sending", func(t *testing.T) {
		// create fake client
		kube := testInternal.NewFakeClient()

		labels := map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}
		if namespaceErr := kube.CreateNamespace(context.TODO(), "test", labels); namespaceErr != nil {
			t.Fatalf("Error injecting namespace add: %v", namespaceErr)
		}

		cfg := structInternal.SchedulerConfig{
			Namespace:   "test",
			Label:       "volume-cleaner/unattached-time",
			GracePeriod: 180,
			TimeFormat:  "2006-01-02_15-04-05Z",
			DryRun:      true,
			NotifTimes:  []int{1, 2, 3, 30},
			EmailCfg: structInternal.EmailConfig{
				BaseURL:         "https://api.notification.canada.ca",
				Endpoint:        "/v2/notifications/email",
				EmailTemplateID: "Random Template",
				APIKey:          "Random APIKEY",
			},
		}

		sort.Slice(cfg.NotifTimes, func(i, j int) bool {
			return cfg.NotifTimes[i] > cfg.NotifTimes[j] // Descending
		})

		type testCase struct {
			timestamp     string
			expectedValue bool
			currNotif     int
		}

		now := time.Now()

		testCases := []testCase{
			{
				timestamp:     now.Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     0,
			},
			{
				timestamp:     now.Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now.Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 150).Format(cfg.TimeFormat),
				expectedValue: true,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 150).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 150).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 150).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 150).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour).Format(cfg.TimeFormat),
				expectedValue: true,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23).Format(cfg.TimeFormat),
				expectedValue: true,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23 - time.Minute*59 - time.Second*59).Format(cfg.TimeFormat),
				expectedValue: true,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23 - time.Minute*59 - time.Second*59).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23 - time.Minute*59 - time.Second*59).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23 - time.Minute*59 - time.Second*59).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23 - time.Minute*59 - time.Second*59).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 151).Format(cfg.TimeFormat),
				expectedValue: true,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 151).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 151).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 151).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 151).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 149).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 149).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 149).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 149).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 149).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 177).Format(cfg.TimeFormat),
				expectedValue: true,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 177).Format(cfg.TimeFormat),
				expectedValue: true,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 177).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 177).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 177).Format(cfg.TimeFormat),
				expectedValue: false,
				currNotif:     4,
			},
		}

		for _, test := range testCases {
			v, err := ShouldSendMail(test.timestamp, test.currNotif, cfg)
			if err != nil {
				t.Fatal("ShouldSendMail failed.")
			}
			assert.Equal(t, v, test.expectedValue)
		}

	})

}
