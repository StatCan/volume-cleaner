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
		log.Printf("Error patching PVC %s from namespace %s: %s", pvc, ns, err)
		return
	}

	log.Printf("Patch successfully applied to %s", pvc)
}

// setting label will add it if doesn't exist
func SetPvcLabel(kube kubernetes.Interface, label string, value string, ns string, pvc string) {
	patchPvcLabel(kube, label, fmt.Sprintf(`"%s"`, value), ns, pvc)
}

// setting label to null (not "null") will remove it
func RemovePvcLabel(kube kubernetes.Interface, label string, ns string, pvc string) {
	pvcObj, err := kube.CoreV1().PersistentVolumeClaims(ns).Get(context.TODO(), pvc, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting PVC to patch: %s", err)
		return
	}
	_, ok := pvcObj.Labels[label]
	if !ok {
		log.Printf("Error removing label %s because label does not exist", label)
		return
	}
	patchPvcLabel(kube, label, "null", ns, pvc)
}
