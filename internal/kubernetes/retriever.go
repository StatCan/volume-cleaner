package kubernetes

import (
	// standard packages
	"context"
	"log"

	// external packages
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	// Internal packages
	structInternal "volume-cleaner/internal/structure"
)

// returns a slice of corev1.Namespace structs

func NsList(kube kubernetes.Interface) []corev1.Namespace {
	ns, err := kube.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/part-of=kubeflow-profile",
	})
	if err != nil {
		// nothing can be done without namespaces so crash the program
		log.Fatalf("Error listing namespaces: %s", err)
	}

	return ns.Items
}

// returns a slice of corev1.PersistentVolumeClaim structs in a given namespace

func PvcList(kube kubernetes.Interface, name string) []corev1.PersistentVolumeClaim {
	pvcs, err := kube.CoreV1().PersistentVolumeClaims(name).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error listing volume claims: %s", err)
	}
	return pvcs.Items
}

// returns a slice of appv1.StatefulSet structs in a given namespace

func StsList(kube kubernetes.Interface, name string) []appv1.StatefulSet {
	sts, err := kube.AppsV1().StatefulSets(name).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error listing stateful sets: %s", err)
	}
	return sts.Items
}

// returns a slice of corev1.PersistentVolumeClaims that are all unattached (not associated with any statefulset)
// from all namespaces

func FindUnattachedPVCs(kube kubernetes.Interface) []corev1.PersistentVolumeClaim {
	/* map each pvc name to its pvc object
	names are used to calculate set differences, but the actual returned
	slice has the objects, so need to keep both
	*/
	pvcObjects := make(map[string]corev1.PersistentVolumeClaim)

	// list of pvc objects to be concated with the pvcs of each namespace
	fullList := make([]corev1.PersistentVolumeClaim, 0)

	log.Print("Scanning namespaces...")

	for _, namespace := range NsList(kube) {
		log.Printf("Found namespace: %s", namespace.Name)
		log.Print("Scanning persistent volume claims...")

		allPVCs := structInternal.NewSet()
		attachedPVCs := structInternal.NewSet()

		// on first pass, add all pvcs to a set

		for _, claim := range PvcList(kube, namespace.Name) {
			// azure disk will have the same name as the volume
			// e.g pvc-11cabba3-59ba-4671-8561-b871e2657fa6

			// claim.Spec.VolumeName will be an empty string if not bound
			log.Printf("PVC: %s, PV: %s", claim.Name, claim.Spec.VolumeName)

			allPVCs.Add(claim.Name)
			pvcObjects[claim.Name] = claim
		}

		log.Print("Scanning stateful sets...")

		// on second pass, add all pvcs attached to sts to a set

		for _, statefulset := range StsList(kube, namespace.Name) {
			log.Printf("Found stateful set: %s", statefulset.Name)

			// Spec.Volumes will find all the attached pvcs, not pvs

			for _, claim := range statefulset.Spec.Template.Spec.Volumes {
				log.Printf("pvc attached to sts: %s", claim.PersistentVolumeClaim.ClaimName)

				attachedPVCs.Add(claim.PersistentVolumeClaim.ClaimName)
			}

		}

		/*
			Use set difference to find all unattached pvcs
			Because this method operates off allPVCs, it means that any non existent
			pvcs from sts will be automatically filterd out
		*/
		unattachedPVCs := allPVCs.Difference(attachedPVCs)

		for pvc := range unattachedPVCs.GetSet() {
			// add unattached pvcs to full list before loop resets

			fullList = append(fullList, pvcObjects[pvc])
		}

		log.Printf("Found %d total volume claims.", allPVCs.Length())
		log.Printf("Found %d unattached volume claims.", unattachedPVCs.Length())

	}

	return fullList
}
