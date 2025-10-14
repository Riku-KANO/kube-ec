package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 環境変数の取得
	productAddr := getEnv("PRODUCT_SERVICE_ADDR", "localhost:50051")
	userAddr := getEnv("USER_SERVICE_ADDR", "localhost:50052")
	orderAddr := getEnv("ORDER_SERVICE_ADDR", "localhost:50053")
	paymentAddr := getEnv("PAYMENT_SERVICE_ADDR", "localhost:50054")
	port := getEnv("PORT", "8080")

	// gRPCクライアントの初期化
	clients, err := NewGRPCClients(productAddr, userAddr, orderAddr, paymentAddr)
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}

	// ハンドラーの初期化
	handler := NewHandler(clients)

	// Ginの初期化
	r := gin.Default()

	// CORS設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Product routes
		products := api.Group("/products")
		{
			products.GET("", handler.ListProducts)
			products.GET("/:id", handler.GetProduct)
		}

		// User routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}

		// Order routes
		orders := api.Group("/orders")
		{
			orders.POST("", handler.CreateOrder)
			orders.GET("/:id", handler.GetOrder)
			orders.GET("", handler.ListOrders)
		}

		// Payment routes
		payments := api.Group("/payments")
		{
			payments.POST("", handler.CreatePayment)
			payments.POST("/process", handler.ProcessPayment)
		}
	}

	log.Printf("Gateway is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
