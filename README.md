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
