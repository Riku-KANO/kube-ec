package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	productpb "github.com/yourusername/kube-ec/proto/product"
	userpb "github.com/yourusername/kube-ec/proto/user"
	orderpb "github.com/yourusername/kube-ec/proto/order"
	paymentpb "github.com/yourusername/kube-ec/proto/payment"
	commonpb "github.com/yourusername/kube-ec/proto/common"
)

type Handler struct {
	clients *GRPCClients
}

func NewHandler(clients *GRPCClients) *Handler {
	return &Handler{clients: clients}
}

// Product Handlers
func (h *Handler) ListProducts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	category := c.Query("category")
	search := c.Query("search")

	var pageInt, pageSizeInt int32
	c.ShouldBindQuery(&struct {
		Page     *int32 `form:"page"`
		PageSize *int32 `form:"page_size"`
	}{&pageInt, &pageSizeInt})

	resp, err := h.clients.ProductClient.ListProducts(c.Request.Context(), &productpb.ListProductsRequest{
		Pagination: &commonpb.Pagination{
			Page:     pageInt,
			PageSize: pageSizeInt,
		},
		Category:    category,
		SearchQuery: search,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.clients.ProductClient.GetProduct(c.Request.Context(), &productpb.GetProductRequest{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// User Handlers
func (h *Handler) Register(c *gin.Context) {
	var req userpb.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.clients.UserClient.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) Login(c *gin.Context) {
	var req userpb.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.clients.UserClient.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Order Handlers
func (h *Handler) CreateOrder(c *gin.Context) {
	var req orderpb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.clients.OrderClient.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.clients.OrderClient.GetOrder(c.Request.Context(), &orderpb.GetOrderRequest{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListOrders(c *gin.Context) {
	userID := c.Query("user_id")

	var pageInt, pageSizeInt int32
	c.ShouldBindQuery(&struct {
		Page     *int32 `form:"page"`
		PageSize *int32 `form:"page_size"`
	}{&pageInt, &pageSizeInt})

	resp, err := h.clients.OrderClient.ListOrders(c.Request.Context(), &orderpb.ListOrdersRequest{
		UserId: userID,
		Pagination: &commonpb.Pagination{
			Page:     pageInt,
			PageSize: pageSizeInt,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Payment Handlers
func (h *Handler) CreatePayment(c *gin.Context) {
	var req paymentpb.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.clients.PaymentClient.CreatePayment(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) ProcessPayment(c *gin.Context) {
	var req paymentpb.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.clients.PaymentClient.ProcessPayment(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
