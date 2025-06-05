package main

import (
	"context"
	"log"
	"maps"
	"slices"

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// returns a slice of corev1.Namespace structs

func nsList(kube kubernetes.Interface) []corev1.Namespace {
	ns, err := kube.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/part-of=kubeflow-profile",
	})
	if err != nil {
		log.Fatalf("Error listing namespaces: %v", err)
	}

	log.Println(len(ns.Items))

	return ns.Items
}

// returns a slice of corev1.PersistentVolumeClaim structs

func pvcList(kube kubernetes.Interface, name string) []corev1.PersistentVolumeClaim {
	pvcs, err := kube.CoreV1().PersistentVolumeClaims(name).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing volume claims: %v", err)
	}
	return pvcs.Items
}

// returns a slice of appv1.StatefulSet structs

func stsList(kube kubernetes.Interface, name string) []appv1.StatefulSet {
	sts, err := kube.AppsV1().StatefulSets(name).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing stateful sets: %v", err)
	}
	return sts.Items
}

func findUnattachedPVCs(kube kubernetes.Interface) []corev1.PersistentVolumeClaim {
	pvcObjects := make(map[string]corev1.PersistentVolumeClaim)

	log.Print("Scanning namespaces...")

	for _, namespace := range nsList(kube) {
		log.Printf("Found namespace: %v", namespace.Name)
		log.Print("Scanning persistent volume claims...")

		allPVCs := NewSet()
		attachedPVCs := NewSet()

		// azure disk will have the same name as the volume
		// e.g pvc-11cabba3-59ba-4671-8561-b871e2657fa6

		for _, claim := range pvcList(kube, namespace.Name) {
			// claim.Spec.VolumeName will be an empty string if not bound
			log.Printf("PVC: %v, PV: %v", claim.Name, claim.Spec.VolumeName)

			allPVCs.Add(claim.Name)
			pvcObjects[claim.Name] = claim
		}

		log.Print("Scanning stateful sets...")

		for _, statefulset := range stsList(kube, namespace.Name) {
			log.Printf("Found stateful set: %v", statefulset.Name)

			for _, claim := range statefulset.Spec.Template.Spec.Volumes {
				attachedPVCs.Add(claim.Name)
			}

		}

		unattachedPVCs := allPVCs.Difference(attachedPVCs)

		log.Printf("Found %d total volume claims.", allPVCs.Length())
		log.Printf("Found %d unattached volume claims.", unattachedPVCs.Length())

	}

	return slices.Collect(maps.Values(pvcObjects))

}
