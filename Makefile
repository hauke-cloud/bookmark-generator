.PHONY: build run docker-build docker-push helm-lint helm-install helm-uninstall clean

# Variables
APP_NAME=bookmark-generator
DOCKER_IMAGE=$(APP_NAME):latest
HELM_RELEASE=$(APP_NAME)
HELM_CHART=./helm/$(APP_NAME)

# Build the Go binary
build:
	go build -o $(APP_NAME) ./cmd/$(APP_NAME)

# Run locally
run:
	go run ./cmd/$(APP_NAME)/main.go

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Push Docker image (customize registry as needed)
docker-push:
	docker push $(DOCKER_IMAGE)

# Lint Helm chart
helm-lint:
	helm lint $(HELM_CHART)

# Install Helm chart
helm-install:
	helm install $(HELM_RELEASE) $(HELM_CHART)

# Upgrade Helm chart
helm-upgrade:
	helm upgrade $(HELM_RELEASE) $(HELM_CHART)

# Uninstall Helm chart
helm-uninstall:
	helm uninstall $(HELM_RELEASE)

# Clean build artifacts
clean:
	rm -f $(APP_NAME)

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...
