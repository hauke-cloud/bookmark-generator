# Deployment Guide

## Quick Start

### 1. Build the Docker Image

```bash
docker build -t your-registry/bookmark-generator:latest .
docker push your-registry/bookmark-generator:latest
```

### 2. Update Helm Values

Edit `helm/bookmark-generator/values.yaml` and update the image repository:

```yaml
image:
  repository: your-registry/bookmark-generator
  tag: "latest"
```

### 3. Install with Helm

```bash
# Create namespace (optional)
kubectl create namespace bookmark-generator

# Install the chart
helm install bookmark-generator ./helm/bookmark-generator -n bookmark-generator

# Or with custom values
helm install bookmark-generator ./helm/bookmark-generator \
  --namespace bookmark-generator \
  --set image.repository=your-registry/bookmark-generator \
  --set image.tag=latest \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=bookmarks.example.com \
  --set ingress.hosts[0].paths[0].path=/ \
  --set ingress.hosts[0].paths[0].pathType=Prefix
```

### 4. Access the Application

#### Using port-forward (for testing):
```bash
kubectl port-forward -n bookmark-generator svc/bookmark-generator 8080:80
# Then visit http://localhost:8080
```

#### Using Ingress:
Visit your configured ingress hostname (e.g., http://bookmarks.example.com)

## Configuration Options

### Image Configuration

```yaml
image:
  repository: bookmark-generator
  pullPolicy: IfNotPresent
  tag: "latest"
```

### Ingress Configuration

```yaml
ingress:
  enabled: true
  className: "nginx"
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
  hosts:
    - host: bookmarks.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: bookmark-generator-tls
      hosts:
        - bookmarks.example.com
```

### Resource Limits

```yaml
resources:
  limits:
    cpu: 200m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 64Mi
```

### Autoscaling

```yaml
autoscaling:
  enabled: true
  minReplicas: 1
  maxReplicas: 3
  targetCPUUtilizationPercentage: 80
```

## Security Best Practices

1. **Service Account**: The application uses a dedicated service account with minimal RBAC permissions (only read access to ingresses).

2. **Security Context**: The pod runs as non-root user (65534) with read-only root filesystem.

3. **Network Policies**: Consider adding network policies to restrict traffic to/from the application.

4. **TLS**: Always use TLS in production. Configure the ingress with TLS certificates.

## Monitoring

### Health Checks

The application provides two endpoints:

- `/health` - Health check (verifies K8s connectivity)
- `/readiness` - Readiness check

These are configured in the deployment with:

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: http
  initialDelaySeconds: 10
  periodSeconds: 30

readinessProbe:
  httpGet:
    path: /readiness
    port: http
  initialDelaySeconds: 5
  periodSeconds: 10
```

### Logs

View application logs:

```bash
kubectl logs -n bookmark-generator -l app.kubernetes.io/name=bookmark-generator -f
```

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -n bookmark-generator
kubectl describe pod -n bookmark-generator <pod-name>
```

### Check RBAC Permissions

```bash
kubectl auth can-i list ingresses --as=system:serviceaccount:bookmark-generator:bookmark-generator -n default
```

### Test Connectivity

```bash
# Port forward and test
kubectl port-forward -n bookmark-generator svc/bookmark-generator 8080:80

# In another terminal
curl http://localhost:8080/health
curl http://localhost:8080/
```

## Upgrading

```bash
# Upgrade the release
helm upgrade bookmark-generator ./helm/bookmark-generator -n bookmark-generator

# Or with new values
helm upgrade bookmark-generator ./helm/bookmark-generator \
  --namespace bookmark-generator \
  --set image.tag=v1.1.0
```

## Uninstalling

```bash
helm uninstall bookmark-generator -n bookmark-generator
```

## Advanced Configuration

### Using with Multiple Clusters

To monitor multiple clusters, you can deploy the application in each cluster or use a centralized approach with appropriate kubeconfig.

### Custom Namespaces

By default, the application lists ingresses from all namespaces. If you want to limit it to specific namespaces, you'll need to modify the RBAC rules to use Role/RoleBinding instead of ClusterRole/ClusterRoleBinding and update the client code.

### Custom Styling

The web interface template is embedded in the binary. To customize it:

1. Edit `cmd/bookmark-generator/templates/index.html`
2. Rebuild the Docker image
3. Deploy the new version
