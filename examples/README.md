# Example Kubernetes Ingress Resources

This directory contains example Kubernetes resources for testing the bookmark generator.

## Apply Examples

```bash
# Apply example ingresses
kubectl apply -f examples/

# View the bookmarks in the UI
kubectl port-forward -n bookmark-generator svc/bookmark-generator 8080:80
# Visit http://localhost:8080

# Clean up
kubectl delete -f examples/
```

## Resources

- `example-ingresses.yaml` - Sample ingress resources with various configurations
