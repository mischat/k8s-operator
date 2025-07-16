# Kubernetes Custom Resource and Controller POC

This project demonstrates how to create a custom resource in Kubernetes and use a custom controller to manage it. Specifically, it shows how to:

- Define a Custom Resource Definition (CRD) for a simple Go program that logs an environment variable
- Use a second Go program to create instances of this custom resource with different environment variable values
- Implement a custom controller that watches for these custom resources and creates actual Kubernetes workloads

## Components

- **[programA.go](./programA.go)** - A simple Go program that logs the `MY_ENV_VAR` environment variable every 10 seconds
- **[programA-crd.yaml](./programA-crd.yaml)** - Custom Resource Definition for ProgramA instances
- **[programB.go](./programB.go)** - Creates ProgramA custom resources programmatically
- **[controller.go](./controller.go)** - Watches for ProgramA custom resources and creates deployments
- **[Dockerfile](./Dockerfile)** - Builds a container image for programA

## Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop) (with Kubernetes enabled) or Minikube
- [Go](https://golang.org/doc/install)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)

## Setup and Demo

### 1. Start Local Kubernetes Cluster

```bash
# Using Minikube
minikube start
eval $(minikube docker-env)

# Or enable Kubernetes in Docker Desktop settings
```

### 2. Apply the Custom Resource Definition

```bash
kubectl apply -f programA-crd.yaml
```

Verify the CRD was created:
```bash
kubectl get crd programas.crds.example.com
```

### 3. Build the Docker Image

```bash
docker build -t programa:latest .
```

### 4. Initialize Go Module and Install Dependencies

```bash
go mod init k8s-operator
go mod tidy
```

### 5. Create Custom Resources with programB

```bash
go run programB.go
```

This creates 3 ProgramA custom resources with different environment variable values:
- `programa-1` with `MY_ENV_VAR=instance-1`
- `programa-2` with `MY_ENV_VAR=instance-2`
- `programa-3` with `MY_ENV_VAR=instance-3`

Verify the custom resources were created:
```bash
kubectl get programas
kubectl get programas -o custom-columns=NAME:.metadata.name,ENV_VALUE:.spec.envVarValue
```

### 6. Run the Controller

```bash
go run controller.go
```

The controller will:
- Process existing ProgramA custom resources
- Create deployments for each custom resource
- Watch for new/modified custom resources

### 7. Verify the Results

Check that pods are running:
```bash
kubectl get pods
```

Check logs from each pod to see different environment variable values:
```bash
kubectl logs -f <pod-name-1>
kubectl logs -f <pod-name-2>
kubectl logs -f <pod-name-3>
```

You should see each pod logging its respective environment variable value every 10 seconds:
- Pod 1: `Environment Variable Value: instance-1`
- Pod 2: `Environment Variable Value: instance-2`
- Pod 3: `Environment Variable Value: instance-3`

## How It Works

1. **Custom Resource Definition**: Defines a new Kubernetes resource type `ProgramA` with an `envVarValue` field
2. **programB**: Uses the Kubernetes client library to create custom resources programmatically
3. **Controller**: Watches for ProgramA custom resources and creates corresponding deployments with the specified environment variables
4. **programA**: Runs in pods created by the controller, logging the environment variable value

## Cleanup

```bash
# Delete deployments
kubectl delete deployment --all

# Delete custom resources
kubectl delete programas --all

# Delete CRD
kubectl delete crd programas.crds.example.com
```

## Next Steps

This POC demonstrates the basic pattern for Kubernetes operators. To extend this:
- Add more fields to the CRD spec
- Implement proper error handling and retries in the controller
- Add status updates to custom resources
- Implement deletion logic in the controller
- Add validation and defaulting webhooks

## Command Reference

Here are all the commands used during this demo for reference:

### Initial Setup
```bash
# Check if kubectl is installed
which kubectl

# Check if docker is installed
which docker

# Check docker info
docker info | grep -i "operating system"

# Install minikube (if needed)
brew install minikube

# Start minikube
minikube start

# Configure docker to use minikube's docker daemon
eval $(minikube docker-env)
```

### Kubernetes Context Management
```bash
# Check available contexts
kubectl config get-contexts

# Check cluster info
kubectl cluster-info

# Check nodes
kubectl get nodes

# Check minikube status
minikube status
```

### Custom Resource Definition (CRD)
```bash
# Apply the CRD
kubectl apply -f programA-crd.yaml

# Verify CRD was created
kubectl get crd programas.crds.example.com

# List all CRDs
kubectl get crd
```

### Docker Image Building
```bash
# Build the Docker image
docker build -t programa:latest .

# List Docker images
docker images | grep programa
```

### Go Module and Dependencies
```bash
# Initialize Go module
go mod init k8s-operator

# Download dependencies
go mod tidy
```

### Custom Resources Management
```bash
# Create custom resources using programB
go run programB.go

# List custom resources
kubectl get programas

# List with custom columns to show environment values
kubectl get programas -o custom-columns=NAME:.metadata.name,ENV_VALUE:.spec.envVarValue

# Get detailed info about a specific custom resource
kubectl get programa programa-1 -o yaml

# Describe a custom resource
kubectl describe programa programa-1
```

### Controller and Deployments
```bash
# Run the controller (with timeout for demo)
timeout 10 go run controller.go

# List deployments
kubectl get deployments

# List pods
kubectl get pods

# Get pods with more details
kubectl get pods -o wide
```

### Pod Logs and Monitoring
```bash
# Get logs from specific pods (replace with actual pod names)
kubectl logs programa-1-deployment-d5b9d9cf9-mlb4n
kubectl logs programa-2-deployment-75ccc65c75-pvn5j
kubectl logs programa-3-deployment-6d85cf7c59-82xfs

# Follow logs in real-time
kubectl logs -f programa-1-deployment-d5b9d9cf9-mlb4n

# Get logs from all pods with a specific label
kubectl logs -l app=programa-1
```

### Inspection and Debugging
```bash
# Describe pod for troubleshooting
kubectl describe pod <pod-name>

# Get deployment details
kubectl describe deployment programa-1-deployment

# Check events
kubectl get events --sort-by=.metadata.creationTimestamp

# Get all resources
kubectl get all
```

### Cleanup Commands
```bash
# Delete all deployments
kubectl delete deployment --all

# Delete specific deployment
kubectl delete deployment programa-1-deployment

# Delete all custom resources
kubectl delete programas --all

# Delete specific custom resource
kubectl delete programa programa-1

# Delete the CRD (this will also delete all custom resources)
kubectl delete crd programas.crds.example.com

# Stop minikube
minikube stop

# Delete minikube cluster
minikube delete
```

### Useful Shortcuts
```bash
# Use shortnames for custom resources
kubectl get pa  # same as kubectl get programas

# Watch resources in real-time
kubectl get pods -w
kubectl get programas -w

# Get resource usage
kubectl top pods
kubectl top nodes
```
