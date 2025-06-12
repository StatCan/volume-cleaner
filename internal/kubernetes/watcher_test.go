package kubernetes

import (
	// External Imports
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	structInternal "volume-cleaner/internal/structure"
)

func TestWatcherLabelling(t *testing.T) {

	t.Run("successful labelling of attached and detached pvcs", func(t *testing.T) {
		// create fake client
		client := fake.NewClientset()

		// create testing namespace
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test",
			Labels: map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}}}
		_, err := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Error injecting namespace add: %v", err)
		}

		names := []string{"pvc1", "pvc2"}

		// create testing pvcs
		for _, name := range names {
			pvc := &corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "test"}}
			_, err := client.CoreV1().PersistentVolumeClaims("test").Create(context.TODO(), pvc, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Error injecting pvc add: %v", err)
			}
		}

		ctx := context.Background()

		cfg := structInternal.Config{
			Namespace:  "test",
			Label:      "volume-cleaner/unattached-time",
			TimeFormat: "2006-01-02_15-04-05Z",
		}

		go WatchSts(ctx, client, cfg)

		time.Sleep(2 * time.Second)

		// no pvc should have labels right now
		pvcs := PvcList(client, "test")

		_, ok := pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		// mock a stateful set attached to a pvc1

		sts := &appv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sts1",
				Namespace: "test",
			},
			Spec: appv1.StatefulSetSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Volumes: []corev1.Volume{
							{
								Name: "pvc1",
								VolumeSource: corev1.VolumeSource{
									PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
										ClaimName: "pvc1",
									},
								},
							},
						},
					},
				},
			},
		}
		_, err = client.AppsV1().StatefulSets("test").Create(context.TODO(), sts, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Error injecting sts add: %v", err)

		}

		time.Sleep(2 * time.Second)

		// should be no change

		pvcs = PvcList(client, "test")

		_, ok = pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		// delete sts

		err = client.AppsV1().StatefulSets("test").Delete(context.TODO(), "sts1", metav1.DeleteOptions{})
		if err != nil {
			t.Fatalf("Error injecting event add: %v", err)
		}

		time.Sleep(2 * time.Second)

		// should have new labels

		pvcs = PvcList(client, "test")

		_, ok = pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, true)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		ctx.Done()

	})
}

func TestInitialScan(t *testing.T) {

	t.Run("successful labelling of unatatched pvcs on controller startup", func(t *testing.T) {
		// create fake client
		client := fake.NewClientset()

		// create testing namespace
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test",
			Labels: map[string]string{"app.kubernetes.io/part-of": "kubeflow-profile"}}}
		_, err := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Error injecting namespace add: %v", err)
		}

		names := []string{"pvc1", "pvc2"}

		// create testing pvcs
		for _, name := range names {
			pvc := &corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "test"}}
			_, err := client.CoreV1().PersistentVolumeClaims("test").Create(context.TODO(), pvc, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Error injecting pvc add: %v", err)
			}
		}

		// no pvc should have labels right now
		pvcs := PvcList(client, "test")

		_, ok := pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, false)

		ctx := context.Background()

		cfg := structInternal.Config{
			Namespace:  "test",
			Label:      "volume-cleaner/unattached-time",
			TimeFormat: "2006-01-02_15-04-05Z",
		}

		InitialScan(client, cfg)

		time.Sleep(2 * time.Second)

		// should have new labels

		pvcs = PvcList(client, "test")

		_, ok = pvcs[0].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, true)

		_, ok = pvcs[1].Labels["volume-cleaner/unattached-time"]
		assert.Equal(t, ok, true)

		ctx.Done()

	})
}
