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
	"kubey/server/internal/api/websocket"
	"kubey/server/internal/config"
	"kubey/server/internal/routes"
	"kubey/server/internal/services"
)

func main() {
	cfg := config.Load()

	// Initialize Kubernetes client
	if err := services.InitKubernetesClient(cfg.KubeConfig); err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	routes.Setup(router, hub)

	srv := &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
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
