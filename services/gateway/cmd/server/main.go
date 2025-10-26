package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appuser "github.com/Riku-KANO/kube-ec/services/gateway/internal/application/user"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/infrastructure/grpc"
	httpserver "github.com/Riku-KANO/kube-ec/services/gateway/internal/presentation/http"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/presentation/http/handler"
)

func main() {
	// Load configuration from environment variables
	config := grpc.ClientConfig{
		AuthServiceAddr: getEnv("AUTH_SERVICE_ADDR", "localhost:50052"),
		UserServiceAddr: getEnv("USER_SERVICE_ADDR", "localhost:50051"),
	}

	// Initialize infrastructure layer (gRPC clients)
	grpcClients, err := grpc.NewClients(config)
	if err != nil {
		log.Fatalf("Failed to create gRPC clients: %v", err)
	}
	defer grpcClients.Close()

	// Initialize repositories
	authRepo := grpc.NewAuthRepository(grpcClients.AuthClient)
	userRepo := grpc.NewUserRepository(grpcClients.UserClient)

	// Initialize application services
	// Use authRepo for authentication operations
	userService := appuser.NewService(authRepo)

	// Initialize presentation layer (HTTP handlers)
	userHandler := handler.NewUserHandler(userService)

	// Setup router
	router := httpserver.SetupRouter(userHandler)

	// Create HTTP server
	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting gateway server on port %s...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
