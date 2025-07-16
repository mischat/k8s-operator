# AWS EKS Setup Notes

## Additional EKS-Specific Permissions

### 1. IAM Roles for Service Accounts (IRSA)
If your controller needs to interact with AWS services, you might need IRSA:

```bash
# Create an IAM role for the service account
eksctl create iamserviceaccount \
  --name controller-service-account \
  --namespace default \
  --cluster your-cluster-name \
  --attach-policy-arn arn:aws:iam::aws:policy/ReadOnlyAccess \
  --approve
```

### 2. Container Registry Permissions
For pushing images to ECR:

```bash
# Get login token for ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com

# Build and tag images for ECR
docker build -t programb:latest -f Dockerfile.programB .
docker tag programb:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/programb:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/programb:latest

docker build -t controller:latest -f Dockerfile.controller .
docker tag controller:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/controller:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/controller:latest
```

### 3. Network Policies (if using Calico/Cilium)
You might need network policies for pod-to-pod communication:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: controller-network-policy
spec:
  podSelector:
    matchLabels:
      app: controller
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: programb
  egress:
  - to: []
    ports:
    - protocol: TCP
      port: 443  # Kubernetes API server
```

### 4. Pod Security Standards
EKS might have pod security policies or standards:

```yaml
apiVersion: v1
kind: Pod
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    fsGroup: 2000
  containers:
  - name: controller
    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities:
        drop:
        - ALL
```

## Deployment Order

1. Apply CRD
2. Apply RBAC configurations
3. Build and push container images
4. Apply deployments

```bash
# Apply in order
kubectl apply -f programA-crd.yaml
kubectl apply -f programB-rbac.yaml
kubectl apply -f controller-rbac.yaml
kubectl apply -f programB-deployment.yaml
kubectl apply -f controller-deployment.yaml
```

## Monitoring and Logging

### CloudWatch Integration
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluent-bit-config
data:
  fluent-bit.conf: |
    [SERVICE]
        Flush         1
        Log_Level     info
        
    [INPUT]
        Name              tail
        Path              /var/log/containers/*controller*.log
        Parser            docker
        Tag               controller.*
        
    [OUTPUT]
        Name              cloudwatch_logs
        Match             controller.*
        log_group_name    /aws/eks/controller
        log_stream_prefix controller-
        region            us-east-1
```

### Prometheus Metrics
Add metrics endpoints to your applications for monitoring with Prometheus/Grafana in EKS.
