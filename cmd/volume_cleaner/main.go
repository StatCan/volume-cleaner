package main

import (
	"context"
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest" // for InClusterConfig and RESTClient
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// watchPods(clientset)
	watchSts(clientset)
}

func watchSts(clientset *kubernetes.Clientset) {
	// Watch StatefulSets in all namespaces
	watcher, err := clientset.AppsV1().StatefulSets("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Watching for StatefulSet events...")

	for event := range watcher.ResultChan() {
		sts, ok := event.Object.(*appsv1.StatefulSet)
		if !ok {
			continue
		}

		switch event.Type {
		case watch.Deleted:
			fmt.Printf("StatefulSet deleted: %s/%s\n", sts.Namespace, sts.Name)
			fmt.Println(sts)
			dumpManifest(sts)
		}
	}

	for event := range watcher.ResultChan() {
		sts, ok := event.Object.(*appsv1.StatefulSet)
		if !ok {
			continue
		}

		switch event.Type {
		case watch.Deleted:
			fmt.Printf("StatefulSet deleted: %s/%s\n", sts.Namespace, sts.Name)
			fmt.Println(sts)
			dumpManifest(sts)
		}
	}

}

func dumpManifest(sts *appsv1.StatefulSet) {
	out, err := json.MarshalIndent(sts, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal STS: %v\n", err)
		return
	}
	fmt.Println("Last known manifest before deletion:")
	fmt.Println(string(out))
}
