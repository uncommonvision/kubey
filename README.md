# Kubey

Kubey is a Kubernetes cluster monitoring and visualization tool with a React/TypeScript frontend and Go backend that connects to Kubernetes clusters via the Kubernetes API.

## Architecture

This is a monorepo with two main components:

- `web/` - React frontend (TypeScript + Vite)
- `api/` - Go backend (Gin framework)

## Development Commands

### Frontend (web/)

```bash
# Install dependencies
npm install

# Development server (localhost only)
npm run dev

# Development server (accessible at http://0.0.0.0:5173)
npm run dev+

# Build for production
npm run build

# Preview production build
npm run preview
```

### Backend (api/)

```bash
# Install dependencies
go mod download

# Run server
go run cmd/api/main.go

# Build binary
go build -o kubey-api ./cmd/api/main.go
```

### Testing

The project uses `gotestsum` for testing to provide clean, readable test output.

**Run all tests:**

```bash
make test
```

This command uses `gotestsum` for clean output and **does not generate a coverage report**.

**Install gotestsum:**

```bash
go install gotest.tools/gotestsum@latest
```

**Run specific test:**

```bash
# For the backend API
cd api && gotestsum --format testname ./internal/handlers/health_test.go

# For the frontend (if tests exist)
cd web && npm test
```

### Other Make Targets

```bash
make api-deps        # Download Go module dependencies
make api-build       # Compile the API binary
make api-run         # Run API in development mode
make api-clean       # Remove the API binary
make api-test        # Run Go tests using gotestsum

make web-deps        # Install bun dependencies
make web-dev         # Start Vite dev server (localhost only)
make web-dev+        # Start Vite dev server on 0.0.0.0
make web-build       # Production build of the UI
make web-preview     # Preview built UI (localhost only)
make web-preview+    # Preview built UI on 0.0.0.0
make web-clean       # Remove generated UI files
make web-test        # Run frontend test suite

make dev             # Run API + web dev servers concurrently
make start           # Build both and serve production UI
make clean           # Remove all generated files
make help            # Show all available targets
```

## API Endpoints

REST API (runs on port 8080):

- `GET /api/clusters` - List all clusters
- `GET /api/clusters/:id` - Get cluster details
- `GET /api/clusters/:id/nodes` - Get cluster nodes
- `GET /api/clusters/:id/pods` - Get all pods
- `GET /api/clusters/:id/services` - Get all services
- `GET /api/clusters/:id/deployments` - Get all deployments
- `GET /api/clusters/:id/namespaces` - Get namespaces with resources
- `GET /health` - Health check

## Environment Configuration

### Backend (.env in api/)

Create `api/.env` (optional, falls back to defaults):

```
ENVIRONMENT=development
HOST=localhost
PORT=8080
LOG_LEVEL=info
KUBECONFIG=  # Leave empty to use default ~/.kube/config
HTTP_READ_TIMEOUT=10  # HTTP read timeout in seconds (default: 10)
HTTP_WRITE_TIMEOUT=10  # HTTP write timeout in seconds (default: 10)
HTTP_IDLE_TIMEOUT=30  # HTTP idle timeout in seconds (default: 30)
```

The backend allows CORS from any localhost or 127.0.0.1 origin on any port for local development flexibility.

## Testing

The backend testing setup has been improved with the following changes:

- **Health endpoint** has been moved to `api/internal/handlers/health.go` in the shared `handlers` package
- **Unit tests** are included in `api/internal/handlers/health_test.go`
- **Test runner** uses `gotestsum` instead of direct `go test` for better output formatting
- **No coverage collection** by default to keep test execution fast

Run tests with:

```bash
make test
```

Or directly with gotestsum:

```bash
cd api && gotestsum --format testname -- ./...
```

## License

See LICENSE file for details.