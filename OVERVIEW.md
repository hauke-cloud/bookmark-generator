# Kubernetes Bookmark Generator - Project Overview

## 📁 Project Structure

```
bookmark-generator/
├── cmd/
│   └── bookmark-generator/
│       ├── main.go                 # Application entry point
│       └── templates/
│           └── index.html          # Web UI template
├── pkg/
│   ├── bookmarks/
│   │   ├── generator.go           # Bookmark generation logic
│   │   └── generator_test.go      # Unit tests
│   └── kubernetes/
│       └── client.go              # K8s client wrapper
├── internal/
│   └── handlers/
│       └── handler.go             # HTTP handlers
├── helm/
│   └── bookmark-generator/
│       ├── Chart.yaml             # Helm chart metadata
│       ├── values.yaml            # Default values
│       ├── values.production.yaml # Production values example
│       └── templates/
│           ├── deployment.yaml
│           ├── service.yaml
│           ├── ingress.yaml
│           ├── serviceaccount.yaml
│           ├── rbac.yaml
│           ├── hpa.yaml
│           ├── _helpers.tpl
│           └── NOTES.txt
├── examples/
│   ├── README.md
│   └── example-ingresses.yaml    # Example K8s ingresses
├── .github/
│   └── workflows/
│       └── build.yaml             # CI/CD pipeline
├── Dockerfile
├── Makefile
├── go.mod
├── README.md
├── DEPLOYMENT.md
└── LICENSE
```

## 🚀 Features

### Core Functionality
- ✅ Connects to Kubernetes API using service account
- ✅ Discovers all ingress routes across all namespaces
- ✅ Generates Firefox-compatible HTML bookmark files
- ✅ Generates Chrome-compatible JSON bookmark files
- ✅ Web interface to view and download bookmarks
- ✅ Detects TLS configuration for https:// URLs

### Technical Features
- ✅ Clean architecture with separated concerns
- ✅ Embedded templates (no runtime dependencies)
- ✅ Health and readiness probes
- ✅ Proper error handling and logging
- ✅ Comprehensive unit tests
- ✅ Go best practices (go fmt, go vet compatible)
- ✅ Efficient resource usage (<128Mi memory)

### Kubernetes Integration
- ✅ Service account with minimal RBAC permissions
- ✅ Helm chart for easy deployment
- ✅ Support for autoscaling (HPA)
- ✅ Ingress configuration
- ✅ Security contexts (non-root, read-only filesystem)
- ✅ Pod anti-affinity for HA

## 🔧 Key Components

### 1. Kubernetes Client (`pkg/kubernetes/client.go`)
- In-cluster authentication
- Lists ingresses from all namespaces
- Extracts host, path, and TLS information
- Constructs proper URLs with scheme detection

### 2. Bookmark Generator (`pkg/bookmarks/generator.go`)
- Firefox HTML format (NETSCAPE-Bookmark-file-1)
- Chrome JSON format (with proper structure)
- HTML escaping for security
- Consistent GUID generation

### 3. HTTP Handlers (`internal/handlers/handler.go`)
- `/` - Web interface
- `/firefox/bookmarks.html` - Download Firefox bookmarks
- `/chrome/bookmarks.json` - Download Chrome bookmarks
- `/health` - Health check with K8s connectivity verification
- `/readiness` - Readiness probe

### 4. Web Interface (`templates/index.html`)
- Modern, responsive design
- Lists all discovered ingresses
- Direct download buttons
- Color-coded namespaces
- Clickable URLs

## 📦 Deployment

### Quick Start
```bash
# Build Docker image
docker build -t your-registry/bookmark-generator:latest .

# Install via Helm
helm install bookmark-generator ./helm/bookmark-generator \
  --set image.repository=your-registry/bookmark-generator \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=bookmarks.example.com
```

### Production Deployment
```bash
helm install bookmark-generator ./helm/bookmark-generator \
  -f helm/bookmark-generator/values.production.yaml \
  --set image.repository=your-registry/bookmark-generator \
  --set image.tag=1.0.0
```

## 🔒 Security

### RBAC Permissions
The service account requires:
- `get`, `list`, `watch` on `ingresses` in `networking.k8s.io` API group
- ClusterRole for access across all namespaces

### Security Context
- Runs as non-root user (UID 65534)
- Read-only root filesystem
- No privilege escalation
- Dropped all capabilities

## 🧪 Testing

```bash
# Run unit tests
go test ./... -v

# Run with coverage
go test ./... -coverprofile=coverage.out

# Format code
go fmt ./...

# Vet code
go vet ./...

# Lint Helm chart
helm lint helm/bookmark-generator
```

## 📊 Resource Usage

**Default:**
- CPU: 100m (request), 200m (limit)
- Memory: 64Mi (request), 128Mi (limit)

**With Autoscaling:**
- Min replicas: 1
- Max replicas: 3
- Target CPU: 80%

## 🌐 Browser Support

### Firefox
- Imports via: Bookmarks → Show All Bookmarks → Import and Backup
- Format: HTML (NETSCAPE-Bookmark-file-1)

### Chrome/Chromium
- Imports via: chrome://bookmarks → Import bookmarks
- Format: JSON

## 🔄 CI/CD

GitHub Actions workflow includes:
- Go tests with race detection
- Code coverage reporting
- golangci-lint
- Docker image building
- Helm chart validation

## 📝 Configuration

### Environment Variables
- `PORT`: HTTP server port (default: 8080)

### Helm Values
See `helm/bookmark-generator/values.yaml` for all options:
- Image configuration
- Resource limits
- Ingress settings
- Autoscaling
- Security contexts
- RBAC

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details

## 🙏 Acknowledgments

- Built with Go and the Kubernetes client-go library
- Uses standard Go project layout
- Follows Kubernetes best practices
