package main

import (
	"context"
	"fmt"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Create Kubernetes clients
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating dynamic client: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	// Define the GroupVersionResource for our custom resource
	gvr := schema.GroupVersionResource{
		Group:    "crds.example.com",
		Version:  "v1",
		Resource: "programas",
	}

	fmt.Println("Starting controller...")
	
	// Watch for ProgramA custom resources
	watchlist := dynamicClient.Resource(gvr).Namespace("default")
	watcher, err := watchlist.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}
	defer watcher.Stop()

	// Process existing resources first
	fmt.Println("Processing existing ProgramA resources...")
	list, err := watchlist.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error listing existing resources: %v", err)
	} else {
		for _, item := range list.Items {
			createDeploymentForProgramA(clientset, &item)
		}
	}

	// Watch for new/changed resources
	fmt.Println("Watching for new ProgramA resources...")
	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Added:
			fmt.Println("New ProgramA resource detected")
			if obj, ok := event.Object.(*unstructured.Unstructured); ok {
				createDeploymentForProgramA(clientset, obj)
			}
		case watch.Modified:
			fmt.Println("ProgramA resource modified")
			if obj, ok := event.Object.(*unstructured.Unstructured); ok {
				createDeploymentForProgramA(clientset, obj)
			}
		case watch.Deleted:
			fmt.Println("ProgramA resource deleted")
			// Could implement deletion logic here
		}
	}
}

func createDeploymentForProgramA(clientset *kubernetes.Clientset, obj *unstructured.Unstructured) {
	name := obj.GetName()
	namespace := obj.GetNamespace()
	
	// Extract the environment variable value from the custom resource spec
	spec, found, err := unstructured.NestedMap(obj.Object, "spec")
	if err != nil || !found {
		log.Printf("Error getting spec from %s: %v", name, err)
		return
	}
	
	envVarValue, found, err := unstructured.NestedString(spec, "envVarValue")
	if err != nil || !found {
		log.Printf("Error getting envVarValue from %s: %v", name, err)
		return
	}

	fmt.Printf("Creating deployment for %s with env value: %s\n", name, envVarValue)

	// Create a deployment with the environment variable
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name + "-deployment",
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "programa",
							Image: "programa:latest",
							ImagePullPolicy: corev1.PullNever, // Use local image
							Env: []corev1.EnvVar{
								{
									Name:  "MY_ENV_VAR",
									Value: envVarValue,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create or update the deployment
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Error creating deployment for %s: %v", name, err)
		return
	}

	fmt.Printf("Created deployment %s\n", result.GetName())
}

func int32Ptr(i int32) *int32 { return &i }
