# 🔖 Kubernetes Bookmark Generator

A production-ready Go application that connects to the Kubernetes API, discovers ingress routes, and generates downloadable bookmark files for Firefox and Chrome browsers.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.19+-326CE5?style=flat&logo=kubernetes)](https://kubernetes.io/)

## ✨ Features

### Core Functionality
- 🔍 Automatically discovers all ingress routes across all namespaces
- 🦊 Generate Firefox-compatible HTML bookmark files
- 🌐 Generate Chrome-compatible JSON bookmark files  
- 🎨 Beautiful, responsive web interface
- 🔒 Secure in-cluster service account authentication
- 🚀 Fast and lightweight (<128Mi memory)
- 🏥 Built-in health and readiness probes

### Technical Highlights
- ✅ Go best practices and clean architecture
- ✅ Comprehensive unit tests
- ✅ Helm chart for easy deployment
- ✅ Kubernetes security best practices
- ✅ TLS/HTTPS detection
- ✅ Production-ready configuration

## 📚 Documentation

- **[Quick Start Guide](QUICKSTART.md)** - Get running in 5 minutes
- **[Deployment Guide](DEPLOYMENT.md)** - Production deployment details
- **[Project Overview](OVERVIEW.md)** - Architecture and design

## 🚀 Quick Start

### Prerequisites

- Kubernetes cluster (v1.19+)
- kubectl configured
- Helm 3.x
- Docker (for building)

### 1. Build and Deploy

```bash
# Build Docker image
docker build -t bookmark-generator:latest .

# Install with Helm
helm install bookmark-generator ./helm/bookmark-generator \
  --create-namespace \
  --namespace bookmark-generator

# Access the application
kubectl port-forward -n bookmark-generator svc/bookmark-generator 8080:80
# Visit http://localhost:8080
```

### 2. Use the Application

1. Open the web interface in your browser
2. View all discovered ingress routes
3. Click "Download for Firefox" or "Download for Chrome"
4. Import the bookmarks into your browser

See [QUICKSTART.md](QUICKSTART.md) for detailed instructions.

## 📖 Usage Examples

### Basic Deployment
```bash
helm install bookmark-generator ./helm/bookmark-generator
```

### With Ingress
```bash
helm install bookmark-generator ./helm/bookmark-generator \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=bookmarks.example.com
```

### Production with Auto-scaling
```bash
helm install bookmark-generator ./helm/bookmark-generator \
  -f helm/bookmark-generator/values.production.yaml
```

## 🔌 API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /` | Web interface |
| `GET /firefox/bookmarks.html` | Download Firefox bookmarks |
| `GET /chrome/bookmarks.json` | Download Chrome bookmarks |
| `GET /health` | Health check (with K8s connectivity test) |
| `GET /readiness` | Readiness probe |

## 🏗️ Architecture

```
┌─────────────┐      ┌──────────────────┐      ┌─────────────┐
│   Browser   │─────▶│  Web Interface   │      │ Kubernetes  │
│             │      │  (Port 8080)     │      │  API Server │
└─────────────┘      └──────────────────┘      └─────────────┘
                            │                           │
                            │      Service Account      │
                            └──────────────────────────▶│
                                                        │
                            ┌───────────────────────────┘
                            │
                     ┌──────▼─────────┐
                     │ Ingress Routes │
                     │ All Namespaces │
                     └────────────────┘
```

## 🔒 Security & RBAC

The application uses minimal RBAC permissions:
- **ClusterRole**: Read-only access to ingresses
- **Permissions**: `get`, `list`, `watch` on `networking.k8s.io/ingresses`
- **Security Context**: Non-root user, read-only filesystem, no privilege escalation

## 🧪 Testing

```bash
# Run tests
make test

# Or manually
go test ./... -v

# With coverage
go test ./... -coverprofile=coverage.out
```

## 🛠️ Development

```bash
# Install dependencies
make deps

# Format code
make fmt

# Vet code  
make vet

# Build locally
make build

# Run locally (requires kubeconfig)
make run
```

## 📊 Project Structure

```
bookmark-generator/
├── cmd/bookmark-generator/    # Application entry point
├── pkg/                        # Public packages
│   ├── bookmarks/             # Bookmark generation
│   └── kubernetes/            # K8s client
├── internal/handlers/         # HTTP handlers
├── helm/bookmark-generator/   # Helm chart
├── examples/                  # Example resources
└── .github/workflows/         # CI/CD
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Built with:
- [Go](https://go.dev/) - Programming language
- [Kubernetes client-go](https://github.com/kubernetes/client-go) - Kubernetes API client
- [Helm](https://helm.sh/) - Kubernetes package manager

## 📮 Support

- 📖 Check the [documentation](QUICKSTART.md)
- 🐛 [Report issues](https://github.com/hauke-cloud/bookmark-generator/issues)
- 💡 [Request features](https://github.com/hauke-cloud/bookmark-generator/issues)

---

Made with ❤️ for the Kubernetes community
