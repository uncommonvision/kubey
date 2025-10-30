# Middleware Architecture

This document describes the middleware layer implemented in the Kubey server.

## Overview

The server uses a modular middleware architecture with four core middleware components that handle cross-cutting concerns for all HTTP requests.

## Middleware Components

### 1. Recovery Middleware (`internal/middlewares/recovery/recover.go`)
- **Purpose**: Catches panics and prevents server crashes
- **Behavior**: 
  - Recovers from any panic during request processing
  - Logs panic details with request ID for debugging
  - Returns a generic "internal server error" JSON response with HTTP 500
  - Allows other middleware to complete gracefully

### 2. Request ID Middleware (`internal/middlewares/request/request_id.go`)
- **Purpose**: Generates and tracks unique request identifiers
- **Behavior**:
  - Checks for existing `X-Request-ID` header in incoming requests
  - Generates a new UUID if none provided
  - Stores request ID in Gin context for other middleware/handlers
  - Adds `X-Request-ID` header to response
  - Enables request tracing across services

### 3. Logging Middleware (`internal/middlewares/logging/logger.go`)
- **Purpose**: Provides structured request logging
- **Behavior**:
  - Logs request details: method, path, status code, latency, client IP
  - Includes request ID for correlation
  - Uses standard Go `log` package for output
  - Logs after request processing completes

### 4. CORS Middleware (`internal/middlewares/security/cors.go`)
- **Purpose**: Handles Cross-Origin Resource Sharing (CORS)
- **Behavior**:
  - Uses `github.com/gin-contrib/cors` for CORS handling
  - Configurable allowed origins via `ALLOWED_ORIGINS` environment variable
  - Exposes `X-Request-ID` header to frontend
  - Allows common HTTP methods and headers for API calls

## Middleware Order

Middleware is registered in the following order in `cmd/api/main.go`:

1. **Recovery** - Must be first to catch panics from all subsequent handlers
2. **Request ID** - Generates trace identifier before logging occurs
3. **Logging** - Captures request details with request ID included
4. **CORS** - Handles preflight requests and adds appropriate headers

## Configuration

### Environment Variables

All middleware configuration is handled through environment variables defined in `.env.example`:

```bash
# CORS configuration
ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173

# Request ID configuration
REQUEST_ID_HEADER=X-Request-ID
```

### Config Structure

The `ApiConfig` struct in `internal/config/config.go` includes:

```go
type ApiConfig struct {
    // ... existing fields ...
    AllowedOrigins   []string // CORS allowed origins
    RequestIDHeader  string   // Header name for request ID
}
```

## Usage

### Adding New Middleware

To add new middleware:

1. Create a new directory under `internal/middlewares/`
2. Implement a function that returns `gin.HandlerFunc`
3. Add the middleware to the registration chain in `cmd/api/main.go`
4. Update configuration if needed in `internal/config/config.go`

### Request Context

All handlers can access the request ID via:

```go
requestID := c.GetString("RequestID")
```

### Example Handler with Request ID

```go
func MyHandler(c *gin.Context) {
    requestID := c.GetString("RequestID")
    log.Printf("Processing request: %s", requestID)
    
    // Handler logic here
    
    c.JSON(200, gin.H{"message": "success"})
}
```

## Development vs Production

- **Development**: Uses default localhost CORS origins
- **Production**: Uses configurable origins from `ALLOWED_ORIGINS` environment variable

## Dependencies

The middleware layer adds the following dependencies:

- `github.com/google/uuid` - For request ID generation
- `github.com/gin-contrib/cors` - For CORS handling

These are automatically managed through Go modules.