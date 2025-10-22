package main

import (
	"log"
	"os"

	appuser "github.com/Riku-KANO/kube-ec/services/gateway/internal/application/user"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/infrastructure/grpc"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/presentation/http"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/presentation/http/handler"
)

func main() {
	// Load configuration from environment variables
	config := grpc.ClientConfig{
		UserServiceAddr: getEnv("USER_SERVICE_ADDR", "localhost:50051"),
	}

	// Initialize infrastructure layer (gRPC clients)
	grpcClients, err := grpc.NewClients(config)
	if err != nil {
		log.Fatalf("Failed to create gRPC clients: %v", err)
	}

	// Initialize repositories
	userRepo := grpc.NewUserRepository(grpcClients.UserClient)

	// Initialize application services
	userService := appuser.NewService(userRepo)

	// Initialize presentation layer (HTTP handlers)
	userHandler := handler.NewUserHandler(userService)

	// Setup router
	router := http.SetupRouter(userHandler)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Starting gateway server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
