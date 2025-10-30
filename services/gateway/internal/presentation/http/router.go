package http

import (
	openapi "github.com/Riku-KANO/kube-ec/services/gateway/internal/api"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/presentation/http/handler"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures HTTP routes
func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Register OpenAPI routes using generated handler wrapper
		openapi.RegisterHandlers(v1, userHandler)
	}

	return r
}
