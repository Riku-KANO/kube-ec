# Helm Charts for kube-ec Services

This directory contains Helm charts for deploying kube-ec microservices to Kubernetes.

## Available Charts

- **kube-ec-auth**: Authentication service (gRPC, port 50052)
- **kube-ec-user**: User management service (gRPC, port 50051)
- **kube-ec-gateway**: API Gateway (HTTP, port 80)
- **kube-ec-web**: Frontend application

## Prerequisites

- Kubernetes cluster
- Helm 3.x installed
- kubectl configured
- Required secrets created (see below)

## Quick Start

### 1. Create Namespace and Secrets

```bash
# Create namespace
kubectl create namespace kube-ec

# Create PostgreSQL credentials
kubectl create secret generic postgres-credentials \
  --namespace=kube-ec \
  --from-literal=username=postgres \
  --from-literal=password=$(openssl rand -base64 32) \
  --from-literal=database=kube_ec

# Get the generated password
POSTGRES_PASSWORD=$(kubectl get secret postgres-credentials -n kube-ec -o jsonpath='{.data.password}' | base64 -d)

# Create database URL secret
kubectl create secret generic db-secret \
  --namespace=kube-ec \
  --from-literal=database-url="postgresql://postgres:${POSTGRES_PASSWORD}@postgres:5432/kube_ec?sslmode=disable"

# Create JWT secret for authentication
kubectl create secret generic jwt-secret \
  --namespace=kube-ec \
  --from-literal=secret=$(openssl rand -base64 64)
```

### 2. Deploy PostgreSQL

```bash
kubectl apply -f ../k8s/postgres.yaml
kubectl wait --for=condition=ready pod -l app=postgres -n kube-ec --timeout=60s
```

### 3. Install Services with Helm

#### Auth Service

```bash
helm install kube-ec-auth ./kube-ec-auth \
  --namespace=kube-ec \
  --set image.repository=kube-ec-auth \
  --set image.tag=latest \
  --set image.pullPolicy=Never  # For local development
```

#### User Service

```bash
helm install kube-ec-user ./kube-ec-user \
  --namespace=kube-ec \
  --set image.repository=kube-ec-user \
  --set image.tag=latest \
  --set image.pullPolicy=Never  # For local development
```

#### Gateway Service

```bash
helm install kube-ec-gateway ./kube-ec-gateway \
  --namespace=kube-ec \
  --set image.repository=kube-ec-gateway \
  --set image.tag=latest \
  --set image.pullPolicy=Never  # For local development
```

## Configuration

### Common Values

All charts support these common configuration options:

```yaml
# Replica count
replicaCount: 2

# Image configuration
image:
  repository: ghcr.io/riku-kano/kube-ec-auth
  pullPolicy: Always
  tag: "latest"

# Resource limits
resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 100m
    memory: 128Mi

# Autoscaling
autoscaling:
  enabled: false
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
```

### Service-Specific Configuration

#### Auth Service (kube-ec-auth)

```yaml
# Environment variables
env:
  grpcPort: "50052"

# Secrets
secrets:
  databaseUrl:
    secretName: db-secret
    key: database-url
  jwtSecret:
    secretName: jwt-secret
    key: secret

# Health checks (gRPC uses TCP probes)
livenessProbe:
  tcpSocket:
    port: grpc
  initialDelaySeconds: 10
  periodSeconds: 10

readinessProbe:
  tcpSocket:
    port: grpc
  initialDelaySeconds: 5
  periodSeconds: 5
```

#### Gateway Service (kube-ec-gateway)

```yaml
# Environment variables
env:
  port: "8080"
  ginMode: "release"
  authServiceAddr: "kube-ec-auth:50052"
  userServiceAddr: "kube-ec-user:50051"

# Health checks (HTTP)
livenessProbe:
  httpGet:
    path: /health
    port: http
  initialDelaySeconds: 10
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /health
    port: http
  initialDelaySeconds: 5
  periodSeconds: 5
```

## Customizing Values

### Using Custom Values File

Create a `values-custom.yaml`:

```yaml
replicaCount: 3

resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 250m
    memory: 256Mi

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
```

Install with custom values:

```bash
helm install kube-ec-auth ./kube-ec-auth \
  --namespace=kube-ec \
  --values values-custom.yaml
```

### Setting Individual Values

```bash
helm install kube-ec-auth ./kube-ec-auth \
  --namespace=kube-ec \
  --set replicaCount=3 \
  --set image.tag=v1.2.3 \
  --set resources.limits.cpu=1000m
```

## Upgrading

```bash
# Upgrade auth service
helm upgrade kube-ec-auth ./kube-ec-auth \
  --namespace=kube-ec \
  --set image.tag=v1.2.4

# Upgrade with custom values
helm upgrade kube-ec-gateway ./kube-ec-gateway \
  --namespace=kube-ec \
  --values values-custom.yaml
```

## Uninstalling

```bash
# Remove specific release
helm uninstall kube-ec-auth --namespace=kube-ec

# Remove all releases
helm uninstall kube-ec-auth kube-ec-user kube-ec-gateway --namespace=kube-ec

# Remove namespace (includes all resources)
kubectl delete namespace kube-ec
```

## Development

### Local Development with Docker Images

When developing locally, build and use local Docker images:

```bash
# Build images
docker build -t kube-ec-auth:latest -f services/auth/Dockerfile .
docker build -t kube-ec-user:latest -f services/user/Dockerfile .
docker build -t kube-ec-gateway:latest -f services/gateway/Dockerfile .

# Install with local images
helm install kube-ec-auth ./kube-ec-auth \
  --namespace=kube-ec \
  --set image.repository=kube-ec-auth \
  --set image.tag=latest \
  --set image.pullPolicy=Never
```

### Validating Charts

```bash
# Lint chart
helm lint ./kube-ec-auth

# Dry run to preview resources
helm install kube-ec-auth ./kube-ec-auth \
  --namespace=kube-ec \
  --dry-run --debug

# Template to see generated YAML
helm template kube-ec-auth ./kube-ec-auth \
  --namespace=kube-ec > rendered.yaml
```

## Troubleshooting

### Check Release Status

```bash
# List installed releases
helm list --namespace=kube-ec

# Get release details
helm status kube-ec-auth --namespace=kube-ec

# View release values
helm get values kube-ec-auth --namespace=kube-ec
```

### Common Issues

#### ImagePullBackOff

```bash
# Check if using correct image pull policy for local images
helm upgrade kube-ec-auth ./kube-ec-auth \
  --namespace=kube-ec \
  --set image.pullPolicy=Never
```

#### Secret Not Found

```bash
# Verify secrets exist
kubectl get secrets -n kube-ec

# Check secret content
kubectl get secret jwt-secret -n kube-ec -o yaml
```

#### Service Connection Errors

```bash
# Check service endpoints
kubectl get svc -n kube-ec
kubectl get endpoints -n kube-ec

# Check pod logs
kubectl logs -l app.kubernetes.io/name=kube-ec-auth -n kube-ec
```

## Production Deployment

### Using Container Registry

1. Build and push images to your registry:

```bash
docker build -t ghcr.io/riku-kano/kube-ec-auth:v1.0.0 -f services/auth/Dockerfile .
docker push ghcr.io/riku-kano/kube-ec-auth:v1.0.0
```

2. Create image pull secret (if using private registry):

```bash
kubectl create secret docker-registry regcred \
  --namespace=kube-ec \
  --docker-server=ghcr.io \
  --docker-username=YOUR_USERNAME \
  --docker-password=YOUR_TOKEN
```

3. Install with registry images:

```bash
helm install kube-ec-auth ./kube-ec-auth \
  --namespace=kube-ec \
  --set image.repository=ghcr.io/riku-kano/kube-ec-auth \
  --set image.tag=v1.0.0 \
  --set image.pullPolicy=Always \
  --set imagePullSecrets[0].name=regcred
```

### Production Values

Create `values-prod.yaml`:

```yaml
replicaCount: 3

image:
  repository: ghcr.io/riku-kano/kube-ec-auth
  pullPolicy: Always
  tag: "v1.0.0"

imagePullSecrets:
  - name: regcred

resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 250m
    memory: 256Mi

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 20
  targetCPUUtilizationPercentage: 70
  targetMemoryUtilizationPercentage: 80

podDisruptionBudget:
  enabled: true
  minAvailable: 2
```

## Additional Resources

- [Kubernetes Deployment Guide](../k8s/README.md)
- [Helm Documentation](https://helm.sh/docs/)
- [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)
