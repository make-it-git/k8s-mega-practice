package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// Initialize a random seed
	rand.Seed(time.Now().UnixNano())

	// Create in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// Create a Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	for {
		// Get the unscheduled pods
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
			FieldSelector: "spec.nodeName==,spec.schedulerName=simple-scheduler",
		})
		if err != nil {
			panic(err.Error())
		}

		// Process each unscheduled pod
		for _, pod := range pods.Items {
			fmt.Printf("Found unscheduled pod: %s/%s\n", pod.Namespace, pod.Name)

			// Get all nodes
			nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}

			if len(nodes.Items) == 0 {
				fmt.Println("No nodes available in the cluster.")
				continue
			}

			// Select a random node
			node := nodes.Items[rand.Intn(len(nodes.Items))]

			// Bind the pod to the selected node
			binding := &v1.Binding{
				ObjectMeta: metav1.ObjectMeta{
					Name:      pod.Name,
					Namespace: pod.Namespace,
				},
				Target: v1.ObjectReference{
					Kind: "Node",
					Name: node.Name,
				},
			}

			err = clientset.CoreV1().Pods(pod.Namespace).Bind(context.TODO(), binding, metav1.CreateOptions{})
			if err != nil {
				fmt.Printf("Failed to bind pod: %v\n", err)
			} else {
				fmt.Printf("Pod %s/%s scheduled to node %s\n", pod.Namespace, pod.Name, node.Name)
			}
		}

		// Sleep before next iteration
		time.Sleep(3 * time.Second)
	}
}

