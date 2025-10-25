package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appuser "github.com/Riku-KANO/kube-ec/services/gateway/internal/application/user"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/domain/errors"
	api "github.com/Riku-KANO/kube-ec/services/gateway/internal/api"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService *appuser.Service
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *appuser.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register implements POST /auth/register
func (h *UserHandler) Register(c *gin.Context) {
	var req api.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Error{Error: err.Error()})
		return
	}

	// Convert OpenAPI request to application DTO
	input := toRegisterInput(req)

	// Call application service
	output, err := h.userService.Register(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	// Convert application DTO to OpenAPI response
	resp := toRegisterResponse(output)
	c.JSON(http.StatusCreated, resp)
}

// Login implements POST /auth/login
func (h *UserHandler) Login(c *gin.Context) {
	var req api.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Error{Error: err.Error()})
		return
	}

	input := toLoginInput(req)

	output, err := h.userService.Login(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	resp := toLoginResponse(output)
	c.JSON(http.StatusOK, resp)
}

// GetUser implements GET /users/{id}
func (h *UserHandler) GetUser(c *gin.Context, id string) {
	output, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	user := toUserResponse(output)
	c.JSON(http.StatusOK, user)
}

// UpdateUser implements PUT /users/{id}
func (h *UserHandler) UpdateUser(c *gin.Context, id string) {
	var req api.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Error{Error: err.Error()})
		return
	}

	input := toUpdateUserInput(req)

	output, err := h.userService.UpdateUser(c.Request.Context(), id, input)
	if err != nil {
		handleError(c, err)
		return
	}

	user := toUserResponse(output)
	c.JSON(http.StatusOK, user)
}

// DeleteUser implements DELETE /users/{id}
func (h *UserHandler) DeleteUser(c *gin.Context, id string) {
	err := h.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, api.DeleteUserResponse{Success: true})
}

// handleError maps domain errors to HTTP responses
func handleError(c *gin.Context, err error) {
	switch err {
	case errors.ErrInvalidInput:
		c.JSON(http.StatusBadRequest, api.Error{Error: err.Error()})
	case errors.ErrUserNotFound:
		c.JSON(http.StatusNotFound, api.Error{Error: err.Error()})
	case errors.ErrUnauthorized:
		c.JSON(http.StatusUnauthorized, api.Error{Error: err.Error()})
	case errors.ErrEmailExists:
		c.JSON(http.StatusConflict, api.Error{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, api.Error{Error: "internal server error"})
	}
}
