package main

import (
	"context"
	"fmt"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Config struct {
	DryRun      bool
	GracePeriod int
}

func main() {
	fmt.Println("Volume cleaner started.")

	cfg := Config{
		DryRun:      os.Getenv("DRY_RUN") == "true",
		GracePeriod: 30,
	}

	kubeClient, err := initKubeClient()
	if err != nil {
		log.Fatalf("Error creating kube client: %v", err)
	}

	cleanVolumes(kubeClient, cfg)

}

// pointers?

func initKubeClient() (*kubernetes.Clientset, error) {
	// service runs inside cluster as a pod, therefore will use in-cluster config
	// to connect with cluster

	cfg, err := rest.InClusterConfig()
	if err == nil {
		return kubernetes.NewForConfig(cfg)
	}
	return nil, err
}

func findUnattachedPVCs(kube kubernetes.Interface) {
	allPVCs := NewSet()
	attachedPVCs := NewSet()
	bindings := make(map[string]string)

	log.Print("Scanning namespaces...")

	ns, err := kube.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/part-of=kubeflow-profile",
	})
	if err != nil {
		log.Fatalf("Error listing namespaces: %v", err)
	}

	for _, namespace := range ns.Items {
		log.Printf("Found Kubeflow namespace: %v", namespace.Name)
		log.Print("Scanning persistent volume claims...")

		allPVCs.Clear()
		attachedPVCs.Clear()

		pvcs, err := kube.CoreV1().PersistentVolumeClaims(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing volume claims: %v", err)
		}

		// azure disk will have the same name as the volume
		// e.g pvc-d88040d...

		for _, claim := range pvcs.Items {
			// claim.Spec.VolumeName will be an empty string if not bound
			log.Printf("PVC: %v, PV: %v", claim.Name, claim.Spec.VolumeName)

			allPVCs.Add(claim.Name)
			bindings[claim.Name] = claim.Spec.VolumeName

		}

		log.Print("Scanning stateful sets...")

		sts, err := kube.AppsV1().StatefulSets(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing stateful sets: %v", err)
		}

		for _, statefulset := range sts.Items {
			log.Printf("Found stateful set: %v", statefulset.Name)

			for _, claim := range statefulset.Spec.Template.Spec.Volumes {
				attachedPVCs.Add(claim.Name)
			}

		}

		unattachedPVCs := allPVCs.Difference(attachedPVCs)

		log.Printf("Found %d total volume claims.", allPVCs.Length())
		log.Printf("Found %d unattached volume claims.", unattachedPVCs.Length())

		// PVCs with no PVs
		orphanedPVCs := 0

		attachedOrphanedPVCs := 0
		unattachedOrphanedPVCs := 0

		for v := range allPVCs.list {
			if bindings[v] == "" {
				orphanedPVCs++
				if !unattachedPVCs.Has(v) {
					attachedOrphanedPVCs++
				} else {
					unattachedOrphanedPVCs++
				}
			}
		}

		log.Printf("Found %d orhpaned PVCs.", orphanedPVCs)
		log.Printf("Found %d attached orhpaned PVCs.", attachedOrphanedPVCs)
		log.Printf("Found %d unattached orhpaned PVCs.", unattachedOrphanedPVCs)

	}

}

func cleanVolumes(kube kubernetes.Interface, cfg Config) {
	findUnattachedPVCs(kube)
}
