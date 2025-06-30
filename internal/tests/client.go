package tests

// Abstract out a Fake Kubernetes Client to be used for testing

import (
	// standard packages
	"context"

	// external packages
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

type FakeClient struct {
	kubernetes.Interface
}

// returns a new instance of FakeClient with an embedded fake Kubernetes clientset.
func NewFakeClient() *FakeClient {
	return &FakeClient{Interface: fake.NewClientset()}
}

// creates a Kubernetes Namespace resource with the given name and labels.
func (f *FakeClient) CreateNamespace(ctx context.Context, name string, labels map[string]string) error {
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name, Labels: labels}}
	_, err := f.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	return err
}

// creates a PersistentVolumeClaim in the specified namespace with the given name.
func (f *FakeClient) CreatePersistentVolumeClaim(ctx context.Context, name string, namespace string) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace}}
	_, err := f.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, pvc, metav1.CreateOptions{})
	return pvc, err
}

// creates a basic StatefulSet resource in the specified namespace with the given name.
func (f *FakeClient) CreateStatefulSet(ctx context.Context, name string, namespace string) error {
	sts := &appv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace}}
	_, err := f.AppsV1().StatefulSets(namespace).Create(ctx, sts, metav1.CreateOptions{})
	return err
}

// creates a StatefulSet that references a PersistentVolumeClaim by name.
func (f *FakeClient) CreateStatefulSetWithPvc(ctx context.Context, stsName string, namespace string, pvcName string) error {

	sts := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      stsName,
			Namespace: namespace,
		},
		Spec: appv1.StatefulSetSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: pvcName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: pvcName,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := f.AppsV1().StatefulSets(namespace).Create(ctx, sts, metav1.CreateOptions{})
	return err
}

// deletes a StatefulSet by name from the specified namespace.
func (f FakeClient) DeleteStatefulSet(ctx context.Context, name string, namespace string) error {
	err := f.AppsV1().StatefulSets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	return err
}
