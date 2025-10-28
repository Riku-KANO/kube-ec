# Kubernetes Deployment Guide

This directory contains Kubernetes manifests for deploying the kube-ec services to a Kubernetes cluster.

**⚠️ WARNING**: These configurations are for **DEVELOPMENT/TESTING ONLY**. For production deployments, see the [Production Considerations](#production-considerations) section below.

## Prerequisites

- Kubernetes cluster (minikube, kind, Docker Desktop, or cloud provider)
- kubectl configured to access your cluster
- OpenSSL (for generating secure passwords)

## Quick Start

### 1. Create Namespace

```bash
kubectl create namespace kube-ec
```

### 2. Create Secrets

#### PostgreSQL Credentials

Create a secret with securely generated credentials:

```bash
kubectl create secret generic postgres-credentials \
  --namespace=kube-ec \
  --from-literal=username=postgres \
  --from-literal=password=$(openssl rand -base64 32) \
  --from-literal=database=kube_ec
```

#### Database URL Secret

The `db-secret` is used by the microservices to connect to PostgreSQL:

```bash
# Get the generated password from the postgres-credentials secret
POSTGRES_PASSWORD=$(kubectl get secret postgres-credentials -n kube-ec -o jsonpath='{.data.password}' | base64 -d)

# Create the database URL secret
kubectl create secret generic db-secret \
  --namespace=kube-ec \
  --from-literal=database-url="postgresql://postgres:${POSTGRES_PASSWORD}@postgres:5432/kube_ec?sslmode=disable"
```

#### JWT Secret

Create a secret for JWT token signing:

```bash
kubectl create secret generic jwt-secret \
  --namespace=kube-ec \
  --from-literal=secret=$(openssl rand -base64 64)
```

### 3. Deploy PostgreSQL

```bash
kubectl apply -f postgres.yaml
```

Wait for PostgreSQL to be ready:

```bash
kubectl wait --for=condition=ready pod -l app=postgres -n kube-ec --timeout=60s
```

### 4. Deploy Services with Helm

See the Helm charts in `deploy/helm/` for deploying the microservices.

## Verifying Deployment

Check that all resources are running:

```bash
kubectl get all -n kube-ec
```

View logs for troubleshooting:

```bash
# PostgreSQL logs
kubectl logs -l app=postgres -n kube-ec

# Auth service logs (if deployed via Helm)
kubectl logs -l app.kubernetes.io/name=kube-ec-auth -n kube-ec
```

## Accessing Services

### Port Forwarding (for local testing)

```bash
# Access gateway service
kubectl port-forward -n kube-ec svc/kube-ec-gateway 8080:80

# Test registration
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'
```

## Cleanup

Remove all resources:

```bash
kubectl delete namespace kube-ec
```

## Production Considerations

⚠️ **DO NOT use these configurations in production without modifications:**

### Database

- **Use a managed database service** (AWS RDS, Google Cloud SQL, Azure Database for PostgreSQL)
- If self-hosting:
  - Use StatefulSets with persistent volumes
  - Implement automated backups
  - Set up replication for high availability
  - Use PostgreSQL operators (e.g., CloudNativePG, Zalando postgres-operator)
  - Enable SSL/TLS connections
  - Implement proper monitoring and alerting

### Secrets Management

- **Never commit secrets to version control**
- Use external secret management:
  - Sealed Secrets
  - External Secrets Operator
  - HashiCorp Vault
  - Cloud provider secret managers (AWS Secrets Manager, GCP Secret Manager, Azure Key Vault)
- Implement secret rotation policies
- Use RBAC to restrict secret access

### Security

- Enable Pod Security Standards
- Use Network Policies to restrict traffic between pods
- Implement RBAC with least privilege principle
- Use private container registries
- Scan container images for vulnerabilities
- Enable audit logging
- Implement TLS for all external communication

### Scalability

- Configure Horizontal Pod Autoscaling (HPA)
- Set appropriate resource requests and limits
- Use Pod Disruption Budgets (PDB)
- Implement health checks and readiness probes
- Use multiple replicas for high availability

### Observability

- Implement centralized logging (ELK stack, Loki, etc.)
- Set up metrics collection (Prometheus)
- Add distributed tracing (Jaeger, OpenTelemetry)
- Configure alerting rules
- Create dashboards for monitoring

## Troubleshooting

### PostgreSQL Connection Issues

```bash
# Check PostgreSQL is running
kubectl get pods -l app=postgres -n kube-ec

# Test database connection
kubectl exec -it deployment/postgres -n kube-ec -- psql -U postgres -d kube_ec -c "SELECT 1;"

# Check secret exists and has correct format
kubectl get secret db-secret -n kube-ec -o jsonpath='{.data.database-url}' | base64 -d
```

### Secret Not Found Errors

Ensure all required secrets are created:

```bash
kubectl get secrets -n kube-ec
```

Expected secrets:
- `postgres-credentials`
- `db-secret`
- `jwt-secret`

## Additional Resources

- [Helm Charts Documentation](../helm/README.md)
- [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)
- [12-Factor App Principles](https://12factor.net/)
