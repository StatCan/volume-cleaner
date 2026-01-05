package kubernetes

import (
	// standard packages
	"context"
	"math"
	"testing"
	"time"

	// external packages
	"github.com/stretchr/testify/assert"

	// internal packages
	structInternal "volume-cleaner/internal/structure"
	testInternal "volume-cleaner/internal/utils"
)

func TestFindStale(t *testing.T) {
	t.Run("successful discovery of stale pvcs", func(t *testing.T) {
		// create fake client
		kube := testInternal.NewFakeClient()

		labels := map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}
		if namespaceErr := kube.CreateNamespace(context.TODO(), "test", labels); namespaceErr != nil {
			t.Fatalf("Error injecting namespace add: %v", namespaceErr)
		}

		names := []string{"pvc1", "pvc2"}

		// inject fake pvcs
		for _, name := range names {
			if _, pvcErr := kube.CreatePersistentVolumeClaim(context.TODO(), name, "test"); pvcErr != nil {
				t.Fatalf("Error injecting pvc add: %v", pvcErr)
			}
		}

		schedulerCfg := structInternal.SchedulerConfig{
			Namespace:   "test",
			TimeLabel:   "volume-cleaner/unattached-time",
			NotifLabel:  "volume-cleaner/notification-count",
			IgnoreLabel: "volume-cleaner/ignore",
			GracePeriod: 0,
			TimeFormat:  "2006-01-02_15-04-05Z",
			DryRun:      true,
			NotifTimes:  []int{10},
		}

		deleted, emailed := FindStale(kube, schedulerCfg)

		// nothing was labelled, so nothing should be deleted
		assert.Equal(t, deleted, 0)
		assert.Equal(t, emailed, 0)

		controllerCfg := structInternal.ControllerConfig{
			Namespace:  "test",
			TimeLabel:  "volume-cleaner/unattached-time",
			NotifLabel: "volume-cleaner/notification-count",
			TimeFormat: "2006-01-02_15-04-05Z",
		}

		InitialScan(kube, controllerCfg)

		time.Sleep(5 * time.Second)

		deleted, emailed = FindStale(kube, schedulerCfg)

		assert.Equal(t, deleted, 2)
		assert.Equal(t, emailed, 0)

		SetPvcLabel(kube, "volume-cleaner/ignore", "true", "test", "pvc1")

		deleted, emailed = FindStale(kube, schedulerCfg)

		// now pvc1 should be skipped
		assert.Equal(t, deleted, 1)
		assert.Equal(t, emailed, 0)

		RemovePvcLabel(kube, "volume-cleaner/ignore", "test", "pvc1")

		schedulerCfg.GracePeriod = 5

		deleted, emailed = FindStale(kube, schedulerCfg)

		assert.Equal(t, deleted, 0)
		assert.Equal(t, emailed, 2)

	})

}

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
				// test one day longer than grace period
				timestamp:     time.Now().Add(-time.Hour * 24 * 181).Format(format),
				format:        format,
				gracePeriod:   180,
				expectedValue: true,
			},
			{
				// test one hour shorter than grace period
				timestamp:     time.Now().Add(-time.Hour*24*180 + time.Hour*23).Format(format),
				format:        format,
				gracePeriod:   180,
				expectedValue: false,
			},
			{
				// test one second shorter than grace period
				timestamp:     time.Now().Add(-time.Hour*24*180 + time.Hour*23 + time.Minute*59 + time.Second*59).Format(format),
				format:        format,
				gracePeriod:   180,
				expectedValue: false,
			},
			{
				// test now
				timestamp:     time.Now().Format(format),
				format:        format,
				gracePeriod:   180,
				expectedValue: false,
			},
			{
				// test now with 0 grace period
				timestamp:     time.Now().Format(format),
				format:        format,
				gracePeriod:   0,
				expectedValue: true,
			},
			{
				// test one hour until grace period
				timestamp:     time.Now().Add(time.Hour).Format(format),
				format:        format,
				gracePeriod:   0,
				expectedValue: false,
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
			TimeLabel:   "volume-cleaner/unattached-time",
			NotifLabel:  "volume-cleaner/notification-count",
			IgnoreLabel: "volume-cleaner/ignore",
			GracePeriod: 180,
			TimeFormat:  "2006-01-02_15-04-05Z",
			DryRun:      true,
			NotifTimes:  []int{30, 3, 2, 1},
			EmailCfg: structInternal.EmailConfig{
				BaseURL:         "https://api.notification.canada.ca",
				Endpoint:        "/v2/notifications/email",
				EmailTemplateID: "Random Template",
				APIKey:          "Random APIKEY",
			},
		}

		type testCase struct {
			timestamp     time.Time
			expectedValue bool
			currNotif     int
		}

		now := time.Now()

		testCases := []testCase{
			{
				timestamp:     now,
				expectedValue: false,
				currNotif:     0,
			},
			{
				timestamp:     now,
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now,
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now,
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now,
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 150),
				expectedValue: true,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 150),
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 150),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 150),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 150),
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour),
				expectedValue: true,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour),
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour),
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23),
				expectedValue: true,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23),
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23),
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23 - time.Minute*59 - time.Second*59),
				expectedValue: true,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23 - time.Minute*59 - time.Second*59),
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23 - time.Minute*59 - time.Second*59),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23 - time.Minute*59 - time.Second*59),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour*24*150 - time.Hour*23 - time.Minute*59 - time.Second*59),
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 151),
				expectedValue: true,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 151),
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 151),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 151),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 151),
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 149),
				expectedValue: false,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 149),
				expectedValue: false,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 149),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 149),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 149),
				expectedValue: false,
				currNotif:     4,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 177),
				expectedValue: true,
				currNotif:     0,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 177),
				expectedValue: true,
				currNotif:     1,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 177),
				expectedValue: false,
				currNotif:     2,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 177),
				expectedValue: false,
				currNotif:     3,
			},
			{
				timestamp:     now.Add(-time.Hour * 24 * 177),
				expectedValue: false,
				currNotif:     4,
			},
		}

		for _, test := range testCases {
			v, daysLeft, err := ShouldSendMail(test.timestamp.Format(cfg.TimeFormat), test.currNotif, cfg)
			if err != nil {
				t.Fatal("ShouldSendMail failed.")
			}

			// test that days left until volume deletion is properly calculated
			diff := float64(cfg.GracePeriod) - time.Since(test.timestamp).Hours()/24
			assert.True(t, math.Abs(daysLeft-diff) < 0.5) // account for differences in computations

			assert.Equal(t, v, test.expectedValue)
		}

	})

}
