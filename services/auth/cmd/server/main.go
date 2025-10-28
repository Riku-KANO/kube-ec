package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	authpb "github.com/Riku-KANO/kube-ec/proto/auth"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	appauth "github.com/Riku-KANO/kube-ec/services/auth/internal/application/auth"
	"github.com/Riku-KANO/kube-ec/services/auth/internal/infrastructure/persistence"
	grpchandler "github.com/Riku-KANO/kube-ec/services/auth/internal/presentation/grpc"
)

func main() {
	// Load configuration from environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(5)                  // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")

	// Initialize infrastructure layer (Repository)
	authRepo := persistence.NewAuthRepository(db)

	// Initialize application layer (Service)
	authService := appauth.NewService(
		authRepo,
		jwtSecret,
		24*time.Hour,  // Access token duration
		720*time.Hour, // Refresh token duration (30 days)
	)

	// Initialize presentation layer (gRPC Handler)
	authHandler := grpchandler.NewAuthHandler(authService)

	// Start gRPC server
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, authHandler)

	// Enable reflection for development (e.g., grpcurl)
	reflection.Register(grpcServer)

	log.Printf("Auth service is running on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
