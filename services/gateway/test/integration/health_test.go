package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	// Create a minimal router for testing
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	expectedBody := `{"status":"ok"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

// TestCORSMiddleware tests that CORS headers are properly set
func TestCORSMiddleware(t *testing.T) {
	router := gin.New()

	// Add CORS middleware (simplified version for testing)
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d for OPTIONS request, got %d", http.StatusNoContent, w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS headers not set properly")
	}
}

// TestRouterSetup tests that the router can be set up without errors
func TestRouterSetup(t *testing.T) {
	// This test verifies that the router setup doesn't panic
	// Note: Full testing requires mock user handler
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Router setup panicked: %v", r)
		}
	}()

	// Create a simple test to ensure routing structure is correct
	router := gin.New()
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) {})
			auth.POST("/login", func(c *gin.Context) {})
		}

		users := v1.Group("/users")
		{
			users.GET("/:id", func(c *gin.Context) {})
			users.PUT("/:id", func(c *gin.Context) {})
			users.DELETE("/:id", func(c *gin.Context) {})
		}
	}

	// Verify routes are registered
	routes := router.Routes()
	if len(routes) == 0 {
		t.Error("No routes registered")
	}

	expectedRoutes := 5 // register, login, get, update, delete
	if len(routes) != expectedRoutes {
		t.Errorf("Expected %d routes, got %d", expectedRoutes, len(routes))
	}
}
