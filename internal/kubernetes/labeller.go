package kubernetes

import (
	// standard packages
	"context"
	"fmt"
	"log"

	// external packages
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

// modifies pvc labels
// requires sufficient rbac permissions
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
		log.Printf("[ERROR] Failed to patch PVC %s from NS %s: %s", pvc, ns, err)
		return
	}

	log.Printf("[INFO] Patch successfully applied to PVC %s from NS %s", pvc, ns)
}

// setting label will add it if doesn't exist
func SetPvcLabel(kube kubernetes.Interface, label string, value string, ns string, pvc string) {
	patchPvcLabel(kube, label, fmt.Sprintf(`"%s"`, value), ns, pvc)
}

// setting label to null (not "null") will remove it
func RemovePvcLabel(kube kubernetes.Interface, label string, ns string, pvc string) {
	patchPvcLabel(kube, label, "null", ns, pvc)
}
