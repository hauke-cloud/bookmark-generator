<!-- llm-readme-management spec=1 commit=7919ae738de7d7208754ecf0946acc72f112932e template=golang model=qwen3.8-27b-q4 digest=b3d4b07f2e19 generated=2026-09-08T12:32:58Z -->
<a href="https://hauke.cloud" target="_blank"><img src="https://img.shields.io/badge/home-hauke.cloud-brightgreen" alt="hauke.cloud" style="display: block;" /></a>
<a href="https://github.com/hauke-cloud" target="_blank"><img src="https://img.shields.io/badge/github-hauke.cloud-blue" alt="hauke.cloud Github Organisation" style="display: block;" /></a>
<a href="https://github.com/hauke-cloud/llm-readme-management" target="_blank"><img src="https://img.shields.io/badge/template-golang-orange" alt="Repository type - golang" style="display: block;" /></a>


# Bookmark Generator


<img src="https://raw.githubusercontent.com/hauke-cloud/.github/main/resources/img/organisation-logo-small.png" alt="hauke.cloud logo" width="109" height="123" align="right">


<llm header hint="Name the Go module path and say whether this is a service, a CLI or a library.">

This is a Go web service (module `github.com/hauke-cloud/bookmark-generator`) that you deploy into a Kubernetes cluster via Helm. It lists every Ingress resource across all namespaces and lets you download them as Firefox HTML or Chrome JSON bookmark files. It is intended for cluster operators who want to import their ingress routes into a browser.

</llm>


## :book: Description

<llm description>

`bookmark-generator` is a small Kubernetes web service that turns every Ingress resource in your cluster into a downloadable browser bookmark file. If you manage a cluster with many ingress routes spread across namespaces and want to import them into Firefox or Chrome/Chromium without copy-pasting URLs one by one, this is the tool for that.

It runs as a Deployment inside the cluster, authenticates via the in-cluster service account, and lists `networking.k8s.io/v1` Ingresses across all namespaces. Each host/path pair becomes one bookmark; the scheme is `https` when the host appears in `spec.tls.hosts`, otherwise `http`. A web UI at `/` shows every discovered route and offers two download endpoints.

- Lists all Ingresses cluster-wide and derives one bookmark per host/path (empty paths default to `/`).
- Serves a web UI displaying host, namespace, and name for each route.
- `/firefox/bookmarks.html` returns a Netscape HTML bookmark file grouped by namespace.
- `/chrome/bookmarks.json` returns Chrome-format JSON grouped by namespace.
- `/health` actively queries the K8s API; `/readiness` returns 200 without a cluster call.

The Helm chart is published to `ghcr.io/hauke-cloud/charts` and the container image to `ghcr.io/hauke-cloud/bookmark-generator`.

</llm>


## :clipboard: Requirements

<llm requirements hint="Give the Go version from the go directive in go.mod. Mention Docker only if the repository actually builds an image.">

- **Go 1.21** — pinned in `go.mod`, the CI workflow, and the Dockerfile builder stage. Required for building, testing, and running the binary.
- **Docker** — needed to build the container image (`make docker-build`).
- **Helm 3** — for linting, packaging, and installing the chart.
- **kubectl** — connected to a reachable Kubernetes cluster.
- **A Kubernetes cluster** where you can create cluster-scoped resources. The chart's default RBAC (`rbac.create: true`) creates a ClusterRole and ClusterRoleBinding granting `get`, `list`, and `watch` on `ingresses` across all namespaces.
- **In-cluster service account** — the application authenticates exclusively via `rest.InClusterConfig()`. There is no kubeconfig or token mode, so `make run` will not work outside a cluster.
- **`helm dependency update`** — run this before `helm lint` or `helm install` when working from the local `helm/bookmark-generator/` directory. The chart declares an `oauth2-proxy` 8.3.0 dependency that must be fetched first; the Makefile targets do not do this for you.

</llm>


## 🚀 Getting started

<llm getting_started hint="Cover go build, go run and go test with the real package paths. If a Makefile or Taskfile exists, prefer its targets over raw go commands.">

1. Clone the repository.

```bash
git clone https://github.com/hauke-cloud/bookmark-generator.git
cd bookmark-generator
```

2. Download and tidy Go dependencies (requires Go 1.21).

```bash
make deps
```

3. Compile the binary.

```bash
make build
```

4. Run the unit tests; no cluster is needed.

```bash
make test
```

5. Fetch the chart's external dependencies before installing.

```bash
helm dependency update helm/bookmark-generator
```

6. Install the chart into a new namespace on your cluster.

```bash
helm install bookmark-generator ./helm/bookmark-generator --create-namespace --namespace bookmark-generator
```

7. Forward the service port to your local machine.

```bash
kubectl port-forward -n bookmark-generator svc/bookmark-generator 8080:80
```

Open <http://localhost:8080> in a browser to see the list of Ingress routes and download Firefox (HTML) or Chrome (JSON) bookmark files. Note that `make run` is available but only works inside a cluster, since the app uses in-cluster service-account authentication exclusively.

</llm>


## :airplane: Usage

<llm usage hint="For a library, show a small import-and-call example using real exported identifiers. For a service or CLI, show how it is started and the flags or subcommands it accepts.">

Once the chart is deployed, the service exposes a web UI at `/` and two bookmark download endpoints. The three things you will most often do are:

**Deploy the chart to a cluster**

```bash
helm dependency update helm/bookmark-generator
helm install bookmark-generator ./helm/bookmark-generator \
  --create-namespace --namespace bookmark-generator
```

To expose the UI through an Ingress:

```bash
helm install bookmark-generator ./helm/bookmark-generator \
  --create-namespace --namespace bookmark-generator \
  --set ingress.enabled=true \
  --set ingress.className=nginx \
  --set ingress.hosts[0].host=bookmark-generator.local
```

**Browse and download bookmarks**

Port-forward the Service and open the UI in a browser:

```bash
kubectl port-forward -n bookmark-generator svc/bookmark-generator 8080:80
```

Then visit `http://localhost:8080` to see every Ingress route in the cluster, grouped by namespace. Download a Firefox bookmark file from `/firefox/bookmarks.html` or a Chrome/Chromium JSON file from `/chrome/bookmarks.json`.

**Test with sample Ingresses**

The `examples/` directory ships five sample Ingress resources across the `default`, `production`, and `monitoring` namespaces:

```bash
kubectl apply -f examples/
# browse the UI, then clean up
kubectl delete -f examples/
```

</llm>


## :wrench: Configuration

<llm configuration hint="Environment variables and CLI flags, taken from the flag definitions or the config struct.">

The application reads a single environment variable. The Helm chart exposes the full set of deployment knobs in `helm/bookmark-generator/values.yaml`; the table below covers the entries most relevant to operators.

**Environment variables**

| Name | Default | Description |
|------|---------|-------------|
| `PORT` | `8080` | HTTP listen port for the web server. |

**Helm values** (see `helm/bookmark-generator/values.yaml` for the complete list)

| Name | Type | Default | Required | Description |
|------|------|---------|----------|-------------|
| `image.repository` | string | `ghcr.io/hauke-cloud/bookmark-generator` | No | Container image to deploy. |
| `image.tag` | string | *(unset; falls back to chart `appVersion`)* | No | Image tag override. |
| `replicaCount` | int | `1` | No | Number of replicas (ignored when `autoscaling.enabled` is true). |
| `service.type` | string | `ClusterIP` | No | Kubernetes Service type. |
| `service.port` | int | `80` | No | Service port (maps to container port 8080). |
| `ingress.enabled` | bool | `false` | No | Create an Ingress resource. |
| `ingress.className` | string | `""` | No | Ingress class name. |
| `resources` | object | limits `200m` / `128Mi`; requests `100m` / `64Mi` | No | CPU and memory limits and requests. |
| `autoscaling.enabled` | bool | `false` | No | Enable the HorizontalPodAutoscaler. |
| `autoscaling.minReplicas` | int | `1` | No | HPA minimum replicas. |
| `autoscaling.maxReplicas` | int | `3` | No | HPA maximum replicas. |
| `rbac.create` | bool | `true` | No | Create a ClusterRole and ClusterRoleBinding (get/list/watch on ingresses). |
| `serviceAccount.create` | bool | `true` | No | Create a ServiceAccount. |
| `oauth2-proxy.enabled` | bool | `false` | No | Gate the `oauth2-proxy` chart dependency (no sub-templates are shipped). |

Remaining values—`env`, `livenessProbe`, `readinessProbe`, `podSecurityContext`, `securityContext`, `nodeSelector`, `tolerations`, `affinity`, `nameOverride`, `fullnameOverride`, `imagePullSecrets`, `podAnnotations`—are standard Helm chart knobs documented in `values.yaml`. An example production override lives in `helm/bookmark-generator/values.production.yaml`.

</llm>


## :hammer: Development

<llm development hint="Include go test, go vet and gofmt only where the CI workflows actually run them.">

Before opening a pull request, make sure your changes pass the checks that CI runs on every push and PR to `main` and `develop`.

**Tests**

```bash
make test
```

CI runs the race detector and collects coverage:

```bash
go mod download
go test -v -race -coverprofile=coverage.out ./...
```

**Lint**

CI runs `golangci-lint` (pinned to `latest`; no `.golangci.yml` in the repository, so default linters apply). Run it locally before pushing:

```bash
golangci-lint run
```

The Makefile also provides `make fmt` (`go fmt ./...`) and `make vet` (`go vet ./...`) for local use.

**Helm chart**

The chart declares an `oauth2-proxy` dependency, so you must fetch it before linting or templating:

```bash
cd helm/bookmark-generator
helm dependency update
helm lint .
helm template bookmark-generator .
```

`make helm-lint` runs `helm lint` but does not fetch dependencies first.

**Build**

CI compiles the binary before the Docker build step:

```bash
go build -v -o bookmark-generator ./cmd/bookmark-generator
```

There are no pre-commit hooks, no conventional-commit title check, and no generated files that must be regenerated and committed.

</llm>


## 📄 License

This Project is licensed under the GNU General Public License v3.0

- see the [LICENSE](LICENSE) file for details.


## :coffee: Contributing

To become a contributor, please check out the [CONTRIBUTING](CONTRIBUTING.md) file.


## :email: Contact

For any inquiries or support requests, please open an issue in this
repository or contact us at [contact@hauke.cloud](mailto:contact@hauke.cloud).
