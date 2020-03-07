package main

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// ----------------------- //
	// Start the client setup
	// ----------------------- //

	// input a kubeconfig file
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	// parse kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create a client to talk with API Server
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// ---------------------------------------------------------- //
	// Start talking to k8s by sending Requests to the API Server //
	// ---------------------------------------------------------- //

	// Example 1
	// Get the Pod objects for all the namespaces
	pods, err := client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// Example 2
	// Test if a pod called "pod-with-two-containers" is running in the namespace "myproject"
	namespace := "myproject"
	pod := "pod-with-two-containers"
	_, err = client.ServerVersion()
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			pod, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
	}
}
