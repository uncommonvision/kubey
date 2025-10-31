package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"kubey/api/internal/config"
	"kubey/api/internal/middlewares/logging"
	"kubey/api/internal/middlewares/recovery"
	"kubey/api/internal/middlewares/request"
	"kubey/api/internal/middlewares/security"
	"kubey/api/internal/routes"
	"kubey/api/internal/services/kubernetes"
)

func main() {
	cfg := config.LoadApi()

	// Initialize Kubernetes client
	if err := kubernetes.InitClient(cfg.KubeConfig); err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Register middleware in correct order
	// 1. Recovery - must be first to catch panics
	router.Use(recovery.Recover())
	// 2. Request ID - generate/track request IDs
	router.Use(request.RequestID())
	// 3. Logging - log requests after request ID is set
	router.Use(logging.Logger())
	// 4. CORS - handle cross-origin requests
	router.Use(security.CORS(cfg))

	routes.Setup(router, cfg)

	srv := &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTPReadTimeout,
		WriteTimeout: cfg.HTTPWriteTimeout,
		IdleTimeout:  cfg.HTTPIdleTimeout,
	}

	go func() {
		log.Printf("Server starting on %s:%s", cfg.Host, cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
