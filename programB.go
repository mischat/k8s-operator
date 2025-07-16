package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	// Create Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating dynamic client: %v", err)
	}

	// Define the GroupVersionResource for our custom resource
	gvr := schema.GroupVersionResource{
		Group:    "crds.example.com",
		Version:  "v1",
		Resource: "programas",
	}

	// Create multiple instances of ProgramA with different environment values
	envValues := []string{"instance-1", "instance-2", "instance-3"}

	for i, envValue := range envValues {
		// Create a ProgramA custom resource
		programA := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "crds.example.com/v1",
				"kind":       "ProgramA",
				"metadata": map[string]interface{}{
					"name":      fmt.Sprintf("programa-%d", i+1),
					"namespace": "default",
				},
				"spec": map[string]interface{}{
					"envVarValue": envValue,
				},
			},
		}

		// Create the custom resource in Kubernetes
		result, err := dynamicClient.Resource(gvr).Namespace("default").Create(
			context.TODO(),
			programA,
			metav1.CreateOptions{},
		)
		if err != nil {
			log.Printf("Error creating ProgramA instance %d: %v", i+1, err)
			continue
		}

		fmt.Printf("Created ProgramA instance: %s with env value: %s\n", 
			result.GetName(), envValue)
		
		// Wait a bit between creations
		time.Sleep(2 * time.Second)
	}

	fmt.Println("All ProgramA instances created successfully!")
}
