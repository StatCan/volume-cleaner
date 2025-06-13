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
	} else {
		log.Printf("Patch successfully applied to %s", pvc)
	}

}

func SetPvcLabel(kube kubernetes.Interface, label string, value string, ns string, pvc string) {
	patchPvcLabel(kube, label, fmt.Sprintf(`"%s"`, value), ns, pvc)
}

func RemovePvcLabel(kube kubernetes.Interface, label string, ns string, pvc string) {
	patchPvcLabel(kube, label, "null", ns, pvc)
}
