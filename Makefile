# ────────────────────────────────────────────────────────────────────────
#  Top‑level Makefile – API (Go) + Web (React/TS) using bun
# ────────────────────────────────────────────────────────────────────────

# ==== VARIABLES ==================================================
GO      ?= go
GOFLAGS ?=
BINARY  ?= kubey-api

# Bun (JavaScript/TypeScript) variables
BUN     ?= bun
BUN_RUN ?= $(BUN) run

# ==== API (backend) targets ======================================
.PHONY: api-deps api-build api-run api-clean api-test

api-deps:
	@cd api && $(GO) mod download

api-build: api-deps
	@cd api && $(GO) build -o bin/$(BINARY) ./cmd/api/main.go

api-run: api-deps
	@cd api && $(GO) run ./cmd/api/main.go

api-clean:
	@rm -f api/bin/$(BINARY)

api-test:
	@cd api && $(GO) test ./... -v

# ==== WEB (frontend) targets =====================================
.PHONY: web-deps web-dev web-dev+ web-build web-preview web-preview+ web-clean web-test

web-deps:
	@cd web && $(BUN) install

# Development server – localhost only (default)
web-dev: web-deps
	@cd web && $(BUN_RUN) dev

# Development server – listen on all interfaces (0.0.0.0)
web-dev+: web-deps
	@cd web && HOST=0.0.0.0 $(BUN_RUN) dev

# Production build
web-build: web-deps
	@cd web && $(BUN_RUN) build

# Preview the production build (localhost only)
web-preview: web-build
	@cd web && $(BUN_RUN) preview

# Preview the production build on all interfaces
web-preview+: web-build
	@cd web && HOST=0.0.0.0 $(BUN_RUN) preview

web-clean:
	@rm -rf web/node_modules web/dist

# Run web test suite (assumes a "test" script in package.json)
web-test:
	@cd web && $(BUN_RUN) test

# ==== Convenience targets =========================================
.PHONY: dev start clean help

# Run both API and web dev servers concurrently
dev: api-run web-dev

# Build both and serve the production UI (good for Docker / prod)
start: api-build web-build
	@cd api && ./bin/$(BINARY) & \
	cd web && $(BUN_RUN) preview

# Clean everything
clean: api-clean web-clean
	@echo "🚀 All build artifacts removed"

# ==== HELP =======================================================
help:
	@echo "Makefile targets:"
	@echo "  api-deps       – download Go module dependencies"
	@echo "  api-build      – compile the API binary"
	@echo "  api-run        – run API in development mode"
	@echo "  api-test       – run Go tests"
	@echo "  web-deps       – install bun dependencies"
	@echo "  web-dev        – start Vite dev server (localhost only)"
	@echo "  web-dev+       – start Vite dev server on 0.0.0.0"
	@echo "  web-build      – production build of the UI"
	@echo "  web-preview    – preview built UI (localhost only)"
	@echo "  web-preview+   – preview built UI on 0.0.0.0"
	@echo "  web-test       – run frontend test suite"
	@echo "  dev            – run API + web dev servers concurrently"
	@echo "  start          – build both and serve production UI"
	@echo "  clean          – remove all generated files"
	@echo "  help           – show this help"
