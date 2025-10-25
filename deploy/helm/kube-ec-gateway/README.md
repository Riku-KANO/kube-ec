# kube-ec-gateway Helm Chart

API Gateway service for kube-ec e-commerce platform.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- Docker registry access (GitHub Container Registry)

## Installation

### 1. Build and Push Docker Image

```bash
# Navigate to gateway service directory
cd services/gateway

# Build Docker image
docker build -t ghcr.io/riku-kano/kube-ec-gateway:latest .

# Login to GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u riku-kano --password-stdin

# Push image
docker push ghcr.io/riku-kano/kube-ec-gateway:latest
```

### 2. Install with Helm

```bash
# From project root
helm install gateway ./deploy/helm/kube-ec-gateway

# Or with custom namespace
kubectl create namespace kube-ec
helm install gateway ./deploy/helm/kube-ec-gateway -n kube-ec
```

### 3. Install with Custom Values

```bash
# Override default values
helm install gateway ./deploy/helm/kube-ec-gateway \
  --set image.tag=v1.0.0 \
  --set replicaCount=5 \
  --set service.type=ClusterIP
```

## Configuration

### Key Configuration Options

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas | `3` |
| `image.repository` | Image repository | `ghcr.io/riku-kano/kube-ec-gateway` |
| `image.tag` | Image tag | `latest` |
| `image.pullPolicy` | Image pull policy | `Always` |
| `service.type` | Kubernetes service type | `LoadBalancer` |
| `service.port` | Service port | `80` |
| `env.userServiceAddr` | User service gRPC address | `user-service:50051` |
| `env.productServiceAddr` | Product service gRPC address | `product-service:50051` |
| `env.orderServiceAddr` | Order service gRPC address | `order-service:50051` |
| `env.paymentServiceAddr` | Payment service gRPC address | `payment-service:50051` |
| `resources.limits.cpu` | CPU limit | `500m` |
| `resources.limits.memory` | Memory limit | `512Mi` |
| `autoscaling.enabled` | Enable HPA | `false` |

### Custom values.yaml

Create a `values-prod.yaml` file:

```yaml
replicaCount: 5

image:
  tag: "v1.0.0"

service:
  type: ClusterIP

resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 200m
    memory: 256Mi

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: api.kube-ec.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: gateway-tls
      hosts:
        - api.kube-ec.com
```

Install with custom values:

```bash
helm install gateway ./deploy/helm/kube-ec-gateway -f values-prod.yaml
```

## Upgrading

```bash
# Upgrade to new version
helm upgrade gateway ./deploy/helm/kube-ec-gateway

# Upgrade with new image tag
helm upgrade gateway ./deploy/helm/kube-ec-gateway \
  --set image.tag=v1.1.0
```

## Uninstalling

```bash
helm uninstall gateway
```

## Health Checks

The gateway exposes a `/health` endpoint for health checks:

- **Liveness Probe**: Checks if the container is alive
  - Path: `/health`
  - Initial Delay: 30s
  - Period: 10s

- **Readiness Probe**: Checks if the container is ready to serve traffic
  - Path: `/health`
  - Initial Delay: 10s
  - Period: 5s

## Dependencies

The gateway requires the following microservices to be running:

- `user-service` (gRPC port 50051)
- `product-service` (gRPC port 50051)
- `order-service` (gRPC port 50051)
- `payment-service` (gRPC port 50051)

## Accessing the Gateway

### LoadBalancer (default)

```bash
# Get external IP
kubectl get svc gateway-kube-ec-gateway

# Access the gateway
curl http://<EXTERNAL-IP>/health
```

### ClusterIP (with port-forward)

```bash
# Port forward
kubectl port-forward svc/gateway-kube-ec-gateway 8080:80

# Access locally
curl http://localhost:8080/health
```

## Monitoring

View logs:

```bash
# View all pods logs
kubectl logs -l app.kubernetes.io/name=kube-ec-gateway

# Follow logs
kubectl logs -f deployment/gateway-kube-ec-gateway
```

Check status:

```bash
# Check deployment status
kubectl get deployment gateway-kube-ec-gateway

# Check pod status
kubectl get pods -l app.kubernetes.io/name=kube-ec-gateway

# Describe pod for details
kubectl describe pod <pod-name>
```

## Troubleshooting

### Image Pull Errors

If you encounter image pull errors, ensure:

1. The image exists in the registry
2. You have proper authentication configured
3. `imagePullSecrets` are set if using private registry

### gRPC Connection Errors

If gateway cannot connect to backend services:

1. Verify service names match in `values.yaml`
2. Check that backend services are running
3. Verify network policies allow communication

### Health Check Failures

If health checks fail:

1. Check gateway logs for errors
2. Verify the `/health` endpoint is implemented
3. Adjust probe timeouts if needed

## Development

### Local Testing with Minikube

```bash
# Start minikube
minikube start

# Build image in minikube
eval $(minikube docker-env)
cd services/gateway
docker build -t ghcr.io/riku-kano/kube-ec-gateway:latest .

# Install chart
helm install gateway ./deploy/helm/kube-ec-gateway \
  --set image.pullPolicy=Never

# Access service
minikube service gateway-kube-ec-gateway
```

## License

Copyright © 2024 kube-ec-team
