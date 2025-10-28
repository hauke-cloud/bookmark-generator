# Quick Start Guide

Get the Kubernetes Bookmark Generator up and running in 5 minutes!

## Prerequisites

- Kubernetes cluster (v1.19+)
- kubectl configured
- Helm 3.x
- Docker (for building images)

## Step 1: Build the Docker Image

```bash
# Clone the repository (if not already)
cd bookmark-generator

# Build the Docker image
docker build -t bookmark-generator:latest .

# Tag for your registry (if using remote registry)
docker tag bookmark-generator:latest your-registry.example.com/bookmark-generator:latest

# Push to registry
docker push your-registry.example.com/bookmark-generator:latest
```

**Note:** If using a local cluster (minikube, kind), you can skip pushing and load directly:
```bash
# For minikube
minikube image load bookmark-generator:latest

# For kind
kind load docker-image bookmark-generator:latest
```

## Step 2: Deploy with Helm

### Option A: Default Configuration (ClusterIP Service)

```bash
helm install bookmark-generator ./helm/bookmark-generator \
  --create-namespace \
  --namespace bookmark-generator
```

Access via port-forward:
```bash
kubectl port-forward -n bookmark-generator svc/bookmark-generator 8080:80
# Visit http://localhost:8080
```

### Option B: With Ingress

```bash
helm install bookmark-generator ./helm/bookmark-generator \
  --create-namespace \
  --namespace bookmark-generator \
  --set ingress.enabled=true \
  --set ingress.className=nginx \
  --set ingress.hosts[0].host=bookmarks.example.com \
  --set ingress.hosts[0].paths[0].path=/ \
  --set ingress.hosts[0].paths[0].pathType=Prefix
```

Access via ingress hostname: http://bookmarks.example.com

### Option C: With TLS

```bash
helm install bookmark-generator ./helm/bookmark-generator \
  --create-namespace \
  --namespace bookmark-generator \
  --set ingress.enabled=true \
  --set ingress.className=nginx \
  --set ingress.hosts[0].host=bookmarks.example.com \
  --set ingress.hosts[0].paths[0].path=/ \
  --set ingress.hosts[0].paths[0].pathType=Prefix \
  --set ingress.tls[0].secretName=bookmarks-tls \
  --set ingress.tls[0].hosts[0]=bookmarks.example.com
```

## Step 3: Verify Deployment

```bash
# Check pods
kubectl get pods -n bookmark-generator

# Check service
kubectl get svc -n bookmark-generator

# Check ingress (if enabled)
kubectl get ingress -n bookmark-generator

# View logs
kubectl logs -n bookmark-generator -l app.kubernetes.io/name=bookmark-generator -f
```

## Step 4: Test the Application

### Using port-forward (if no ingress)
```bash
kubectl port-forward -n bookmark-generator svc/bookmark-generator 8080:80
```

### Check health
```bash
curl http://localhost:8080/health
# Should return: OK
```

### View web interface
Open http://localhost:8080 (or your ingress URL) in a browser.

### Download bookmarks
- Click "Download for Firefox" to get HTML bookmarks
- Click "Download for Chrome" to get JSON bookmarks

## Step 5: Import Bookmarks

### Firefox
1. Download the HTML file
2. Open Firefox
3. Press Ctrl+Shift+O (or Cmd+Shift+O on Mac)
4. Click "Import and Backup" → "Import Bookmarks from HTML"
5. Select the downloaded file

### Chrome
1. Download the JSON file
2. Open Chrome
3. Go to chrome://bookmarks
4. Click the three dots (⋮) → "Import bookmarks"
5. Select the downloaded file

## Troubleshooting

### Pods not starting
```bash
# Check pod status
kubectl describe pod -n bookmark-generator <pod-name>

# Check events
kubectl get events -n bookmark-generator --sort-by='.lastTimestamp'
```

### RBAC issues
```bash
# Verify ClusterRole
kubectl get clusterrole bookmark-generator

# Verify ClusterRoleBinding
kubectl get clusterrolebinding bookmark-generator

# Test permissions
kubectl auth can-i list ingresses \
  --as=system:serviceaccount:bookmark-generator:bookmark-generator
```

### No ingresses showing up
```bash
# Check if there are any ingresses in the cluster
kubectl get ingresses --all-namespaces

# Check logs for errors
kubectl logs -n bookmark-generator -l app.kubernetes.io/name=bookmark-generator
```

### Connection issues
```bash
# Test from inside a pod
kubectl run -it --rm debug --image=alpine --restart=Never -- sh
# apk add curl
# curl http://bookmark-generator.bookmark-generator.svc.cluster.local/health
```

## Testing with Example Ingresses

If your cluster doesn't have any ingresses yet:

```bash
# Apply example ingresses
kubectl apply -f examples/example-ingresses.yaml

# Verify they were created
kubectl get ingresses --all-namespaces

# Refresh the bookmark generator UI
# You should now see the example ingresses
```

Clean up examples:
```bash
kubectl delete -f examples/example-ingresses.yaml
```

## Updating

```bash
# Pull latest changes
git pull

# Rebuild Docker image
docker build -t bookmark-generator:latest .
docker push your-registry.example.com/bookmark-generator:latest

# Upgrade Helm release
helm upgrade bookmark-generator ./helm/bookmark-generator \
  --namespace bookmark-generator \
  --set image.tag=latest
```

## Uninstalling

```bash
# Uninstall Helm release
helm uninstall bookmark-generator -n bookmark-generator

# Delete namespace (optional)
kubectl delete namespace bookmark-generator
```

## Next Steps

- Read [DEPLOYMENT.md](DEPLOYMENT.md) for production deployment
- Check [README.md](README.md) for detailed documentation
- See [OVERVIEW.md](OVERVIEW.md) for architecture details

## Need Help?

- Check logs: `kubectl logs -n bookmark-generator -l app.kubernetes.io/name=bookmark-generator`
- Verify RBAC: `kubectl get clusterrole,clusterrolebinding | grep bookmark`
- Test connectivity: `kubectl port-forward -n bookmark-generator svc/bookmark-generator 8080:80`

Happy bookmarking! 🔖
