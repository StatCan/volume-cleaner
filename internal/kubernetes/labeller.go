package kubernetes

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

func patchPvcLabel(kube kubernetes.Interface, label string, value string, ns string, pvc string) {
	patch := []byte(fmt.Sprintf(`{"metadata":{"labels":{"%s":%s}}}`, label, value))
	_, err := kube.CoreV1().PersistentVolumeClaims(ns).Patch(
		context.TODO(),
		pvc,
		types.MergePatchType,
		patch,
		metav1.PatchOptions{},
	)
	if err != nil {
		log.Fatalf("Error patching PVC %s from namespace %s: %v", pvc, ns, err)
	}
	log.Printf("Patch successfully applied to %s", pvc)
}

func SetPvcLabel(kube kubernetes.Interface, label string, value string, ns string, pvc string) {
	patchPvcLabel(kube, label, fmt.Sprintf(`"%s"`, value), ns, pvc)
}

func RemovePvcLabel(kube kubernetes.Interface, label string, ns string, pvc string) {
	patchPvcLabel(kube, label, "null", ns, pvc)
}

// function has no unit tests, unused but might be later in the future
// refer to #67

// func GetPvcLabel(kube kubernetes.Interface, label string, ns string, pvc string) (string, error) {
// 	obj, err := kube.CoreV1().PersistentVolumeClaims(ns).Get(context.TODO(), pvc, metav1.GetOptions{})
// 	if err != nil {
// 		return "", err
// 	}
// 	value, ok := obj.GetLabels()[label]
// 	if !ok {
// 		return "", errors.New("Label does not exist")
// 	}
// 	return value, nil
// }
