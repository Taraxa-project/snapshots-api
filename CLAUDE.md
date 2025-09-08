# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Running
- `make build` - Build the Go application to bin/snapshots-api
- `make run` - Run the application locally (requires Go 1.24.3+)
- `go run ./cmd/server` - Alternative way to run locally

### Testing
- `make test` - Run all tests with race detection and coverage
- `make test-coverage` - Generate HTML coverage report (coverage.html)
- `make dev-test` - Run format, vet, and test (recommended before commits)
- `make ci-test` - Run full CI-like tests including Helm linting

### Code Quality
- `make fmt` - Format Go code
- `make vet` - Run go vet for static analysis
- `make deps` - Download and tidy dependencies

### API Testing (requires running server)
- `make test-api` - Test all API endpoints against localhost:8080

### Docker
- `make docker-build` - Build Docker image as snapshots-api:latest
- `make docker-run` - Run container with port 8080 exposed

### Helm
- `make helm-lint` - Lint the Helm chart
- `make helm-template` - Render Helm templates for validation
- `make helm-install` - Install chart to current kubectl context

## Architecture

### Core Components

**Service Layer** (`internal/service/`):
- `SnapshotService` - Main business logic with 5-minute caching
- Fetches snapshot data from GCP Storage API
- Groups and sorts snapshots by network (mainnet/testnet/devnet) and type (full/light)

**API Layer** (`internal/api/`):
- `Handler` - HTTP handlers with dependency injection pattern
- Routes: `/` (snapshots), `/health`, `/ready`
- Uses standard library's `http.ServeMux`
- Integrates authentication middleware for request filtering

**Parser** (`internal/parser/`):
- `SnapshotParser` - Extracts metadata from snapshot filenames
- Filename format: `{network}-{type}-db-block-{blocknum}-{timestamp}.tar.gz`
- Validates network and type parameters

**Models** (`internal/models/`):
- `Snapshot` - Core data structure with network, type, block number, timestamp
- `NetworkSnapshots` - Response format with full/light snapshot info
- Constants for networks (mainnet/testnet/devnet) and types (full/light)

**Authentication** (`internal/auth/`):
- `Middleware` - Bearer token authentication with standard HTTP headers
- API key extraction and validation
- Request filtering based on authentication status

**Configuration** (`internal/config/`):
- Environment-based config with sensible defaults
- Key variables: `PORT`, `GCP_BUCKET_NAME`, `GCP_BUCKET_URL`, `API_KEYS`
- API key validation against comma-separated environment variable

### Data Flow
1. HTTP request with `?network=` parameter and optional `Authorization` header
2. Authentication middleware checks Bearer token validity
3. Service checks 5-minute cache
4. On cache miss: fetch from `https://storage.googleapis.com/storage/v1/b/{bucket}/o`
5. Parser extracts metadata from filenames using regex
6. Service finds latest snapshots by block number (then timestamp)
7. Service filters response based on authentication (full + light vs light only)
8. JSON response with 5-minute cache headers

### Testing Strategy
- All packages have comprehensive unit tests (`*_test.go`)
- Mocked external dependencies (GCP API calls)
- Integration tests for HTTP handlers
- Coverage reporting enabled

### Key Design Patterns
- **Dependency Injection**: Handler accepts service interface
- **Interface Segregation**: `SnapshotServiceInterface` for testability
- **Caching**: Thread-safe cache with TTL and mutex protection
- **Graceful Shutdown**: 10-second timeout for HTTP server shutdown
- **Error Wrapping**: Context-aware error messages with fmt.Errorf

## Authentication System

**API Key Authentication** (Bearer Token):
- **Light snapshots**: Available without authentication
- **Full snapshots**: Require valid API key in `Authorization: Bearer <token>` header
- Uses standard HTTP Bearer authentication
- Returns 401 Unauthorized with proper WWW-Authenticate header for invalid keys
- Unauthenticated requests receive filtered response (light snapshots only)

**Authentication Components**:
- `auth.Middleware` - Handles token extraction and validation
- `config.IsValidAPIKey()` - Validates API keys against comma-separated list
- Service layer filters snapshots based on authentication status

## Environment Configuration

Default values:
- `PORT=8080`
- `GCP_BUCKET_NAME=taraxa-snapshot`
- `GCP_BUCKET_URL=https://storage.googleapis.com/storage/v1/b/taraxa-snapshot/o`
- `API_KEYS` - Comma-separated list of valid API keys (optional)

## Snapshot Filename Convention

Format: `{network}-{type}-db-block-{blocknumber}-{timestamp}.tar.gz`
- Networks: `mainnet`, `testnet`, `devnet`
- Types: `full`, `light`
- Timestamp: `YYYYMMDD-HHMMSS`

Example: `mainnet-full-db-block-19547931-20250706-062734.tar.gz`