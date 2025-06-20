package kubernetes

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	structInternal "volume-cleaner/internal/structure"
)

func TestIsStale(t *testing.T) {

	t.Run("test successful determination of stale pvcs", func(t *testing.T) {
		format := "2006-01-02_15-04-05Z"

		assert.Equal(t, IsStale(time.Now().Add(-time.Hour*24*180).Format(format), format, 180), false)
		assert.Equal(t, IsStale(time.Now().Add(-time.Hour*24*181).Format(format), format, 180), true)
		assert.Equal(t, IsStale(time.Now().Add(-time.Hour*24*180-time.Hour*23).Format(format), format, 180), false)
		assert.Equal(t, IsStale(time.Now().Add(-time.Hour*24*180-time.Hour*23-time.Minute*59-time.Second*59).Format(format), format, 180), false)
		assert.Equal(t, IsStale(time.Now().Add(-time.Hour*24*1000).Format(format), format, 180), true)
		assert.Equal(t, IsStale(time.Now().Format(format), format, 180), false)
		assert.Equal(t, IsStale(time.Now().Format(format), format, 0), false)
		assert.Equal(t, IsStale(time.Now().Add(-time.Second).Format(format), format, 0), false)
		assert.Equal(t, IsStale(time.Now().Add(-time.Hour*24).Format(format), format, 0), true)

	})
}

func TestShouldSendMail(t *testing.T) {

	t.Run("test successful determination of email sending", func(t *testing.T) {
		// create fake client
		client := fake.NewClientset()

		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test",
			Labels: map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}}}
		_, namespaceErr := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
		if namespaceErr != nil {
			t.Fatalf("Error injecting namespace add: %v", namespaceErr)
		}

		pvc := &corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: "pvc", Namespace: "test"}}
		_, pvcErr := client.CoreV1().PersistentVolumeClaims("test").Create(context.TODO(), pvc, metav1.CreateOptions{})
		if pvcErr != nil {
			t.Fatalf("Error injecting pvc add: %v", pvcErr)
		}

		cfg := structInternal.SchedulerConfig{
			Namespace:   "test",
			Label:       "volume-cleaner/unattached-time",
			GracePeriod: 180,
			TimeFormat:  "2006-01-02_15-04-05Z",
			DryRun:      true,
			NotifTimes:  []int{1, 2, 3, 30},
		}

		assert.Equal(t, ShouldSendMail(time.Now().Format(cfg.TimeFormat), *pvc, cfg), false)
		assert.Equal(t, ShouldSendMail(time.Now().Add(-time.Hour*24*150).Format(cfg.TimeFormat), *pvc, cfg), true)
		assert.Equal(t, ShouldSendMail(time.Now().Add(-time.Hour*24*150-time.Hour).Format(cfg.TimeFormat), *pvc, cfg), true)
		assert.Equal(t, ShouldSendMail(time.Now().Add(-time.Hour*24*150-time.Hour*23).Format(cfg.TimeFormat), *pvc, cfg), true)
		assert.Equal(t, ShouldSendMail(time.Now().Add(-time.Hour*24*150-time.Hour*23-time.Minute*59-time.Second*59).Format(cfg.TimeFormat), *pvc, cfg), true)
		assert.Equal(t, ShouldSendMail(time.Now().Add(-time.Hour*24*151).Format(cfg.TimeFormat), *pvc, cfg), false)
		assert.Equal(t, ShouldSendMail(time.Now().Add(-time.Hour*24*149).Format(cfg.TimeFormat), *pvc, cfg), false)
		assert.Equal(t, ShouldSendMail(time.Now().Add(-time.Hour*24*177).Format(cfg.TimeFormat), *pvc, cfg), true)

	})

}
