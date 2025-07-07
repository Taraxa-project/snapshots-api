# Taraxa Snapshots API

A REST API service that provides information about Taraxa blockchain snapshots stored in Google Cloud Storage. The service fetches snapshot data from a GCP bucket and returns the latest snapshots for different networks (mainnet, testnet, devnet) and types (full, light).

## Features

- **Multi-network support**: mainnet, testnet, devnet
- **Snapshot types**: full and light snapshots
- **Latest snapshot selection**: Returns the most recent snapshot based on block number and timestamp
- **Health and readiness probes**: Kubernetes-ready endpoints
- **Caching**: Built-in caching to reduce API calls to GCP
- **Docker ready**: Multi-stage Docker build with security best practices
- **Kubernetes ready**: Complete Helm chart for deployment
- **CI/CD pipeline**: GitHub Actions with automated testing, building, and security scanning

## API Endpoints

### Get Snapshots
```
GET /?network={network}
```

**Parameters:**
- `network` (required): The blockchain network (`mainnet`, `testnet`, or `devnet`)

**Response:**
```json
{
  "full": {
    "block": 19547931,
    "timestamp": "2025-07-06 06:27",
    "url": "https://storage.googleapis.com/taraxa-snapshot/mainnet-full-db-block-19547931-20250706-062734.tar.gz"
  },
  "light": {
    "block": 19546050,
    "timestamp": "2025-07-06 04:58",
    "url": "https://storage.googleapis.com/taraxa-snapshot/mainnet-light-db-block-19546050-20250706-045815.tar.gz"
  }
}
```

### Health Check
```
GET /health
```

Returns the health status of the service.

### Readiness Check
```
GET /ready
```

Verifies the service can connect to the GCP bucket and is ready to serve requests.

## Configuration

The application can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `GCP_BUCKET_NAME` | `taraxa-snapshot` | GCP bucket name |
| `GCP_BUCKET_URL` | `https://storage.googleapis.com/storage/v1/b/taraxa-snapshot/o` | GCP bucket API URL |

## Development

### Prerequisites

- Go 1.24.3 or later
- Docker (optional)
- Kubernetes cluster with Helm (for deployment)

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/taraxa/snapshots-api.git
   cd snapshots-api
   ```

2. **Install dependencies**
   ```bash
   make deps
   ```

3. **Run tests**
   ```bash
   make test
   ```

4. **Run the application**
   ```bash
   make run
   ```

5. **Test the API**
   ```bash
   # In another terminal
   make test-api
   ```

### Available Make Targets

Run `make help` to see all available targets:

```bash
Available targets:
  help            Show this help message
  build           Build the Go application
  test            Run all tests
  test-coverage   Generate test coverage report
  run             Run the application locally
  clean           Clean build artifacts
  fmt             Format Go code
  vet             Run go vet
  deps            Download and tidy dependencies
  docker-build    Build Docker image
  docker-run      Run Docker container
  helm-lint       Lint Helm chart
  helm-template   Render Helm templates
  helm-install    Install Helm chart (requires kubectl context)
  helm-upgrade    Upgrade Helm chart (requires kubectl context)
  helm-uninstall  Uninstall Helm chart (requires kubectl context)
  dev-setup       Setup development environment
  dev-test        Run development tests (format, vet, test)
  ci-test         Run CI-like tests locally
  test-api        Test API endpoints (requires running server)
```

## Docker

### Build Docker Image
```bash
make docker-build
```

### Run with Docker
```bash
make docker-run
```

### Environment Variables for Docker
```bash
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e GCP_BUCKET_NAME=taraxa-snapshot \
  -e GCP_BUCKET_URL=https://storage.googleapis.com/storage/v1/b/taraxa-snapshot/o \
  snapshots-api:latest
```

## Kubernetes Deployment

### Using Helm

1. **Install the chart**
   ```bash
   helm install snapshots-api ./charts/snapshots-api
   ```

2. **Install with custom values**
   ```bash
   helm install snapshots-api ./charts/snapshots-api \
     --set replicaCount=3 \
     --set resources.requests.memory=256Mi \
     --set ingress.enabled=true \
     --set ingress.hosts[0].host=snapshot.taraxa.io
   ```

3. **Upgrade the deployment**
   ```bash
   helm upgrade snapshots-api ./charts/snapshots-api
   ```

### Configuration Values

Key Helm values you can override:

```yaml
# Number of replicas
replicaCount: 2

# Docker image configuration
image:
  repository: ghcr.io/taraxa/snapshots-api
  tag: "latest"

# Resource limits
resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 100m
    memory: 128Mi

# Ingress configuration
ingress:
  enabled: true
  hosts:
    - host: snapshot.taraxa.io
      paths:
        - path: /
          pathType: Prefix

# Environment variables
env:
  GCP_BUCKET_NAME: "taraxa-snapshot"
  GCP_BUCKET_URL: "https://storage.googleapis.com/storage/v1/b/taraxa-snapshot/o"
```

## CI/CD Pipeline

The project includes a comprehensive GitHub Actions pipeline that:

### On Pull Requests and Pushes to Main:
- Runs all tests with coverage reporting
- Performs code quality checks (formatting, vetting, static analysis)
- Lints Helm charts
- Builds and pushes Docker images to GitHub Container Registry
- Runs security scans with Trivy
- Generates Software Bill of Materials (SBOM)

### On Tag Pushes (v*):
- Creates GitHub releases with deployment instructions
- Tags Docker images with version numbers

### Security Features:
- Vulnerability scanning with Trivy
- SARIF upload to GitHub Security tab
- SBOM generation for supply chain security
- Multi-architecture builds (amd64, arm64)

## Architecture

### Project Structure
```
.
├── cmd/server/           # Application entrypoint
├── internal/
│   ├── api/             # HTTP handlers and routing
│   ├── config/          # Configuration management
│   ├── models/          # Data models
│   ├── parser/          # Snapshot filename parsing
│   └── service/         # Business logic
├── charts/snapshots-api/ # Helm chart
├── .github/workflows/   # CI/CD pipelines
├── Dockerfile           # Container definition
└── Makefile            # Development tasks
```

### Data Flow
1. API receives request with network parameter
2. Service checks cache for recent data
3. If cache miss, fetches snapshot list from GCP bucket
4. Parser extracts metadata from snapshot filenames
5. Service identifies latest snapshots by block number/timestamp
6. Response formatted and returned with caching headers

## Testing

The project includes comprehensive test coverage:

- **Unit tests** for all packages
- **Integration tests** for API endpoints
- **Mock services** for external dependencies
- **Helm chart testing** with template validation

Run tests with coverage:
```bash
make test-coverage
```

## Monitoring and Observability

### Health Checks
- `/health` - Basic service health (always returns 200 if service is running)
- `/ready` - Readiness check (verifies GCP bucket connectivity)

### Logging
- Structured logging for all requests and errors
- Request/response logging with status codes
- Error logging with context

### Metrics
Ready for integration with Prometheus/metrics collection:
- HTTP request duration and status codes (via ingress controller)
- Pod resource usage (via Kubernetes metrics)
- Custom application metrics can be added

## Production Considerations

### Security
- Container runs as non-root user
- Read-only root filesystem
- Security context restrictions
- Regular dependency updates via Dependabot
- Vulnerability scanning in CI

### Performance
- HTTP caching headers set (5-minute cache)
- Internal caching (5-minute TTL)
- Efficient snapshot parsing and selection
- Horizontal pod autoscaling support

### Reliability
- Health and readiness probes
- Graceful shutdown handling
- Resource limits and requests defined
- Multiple replica support

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run `make dev-test` to verify all checks pass
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions and support:
- GitHub Issues: [Create an issue](https://github.com/taraxa/snapshots-api/issues)
- Email: team@taraxa.io