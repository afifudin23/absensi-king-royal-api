package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var payload request.AuthRegisterRequest
	payload.Normalize()

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}

	user, err := h.Service.Register(payload)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyRegistered) {
			common.ErrorHandler(c, common.BadRequestError("Email is already registered"))
			return
		}
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var payload request.AuthLoginRequest
	payload.Normalize()

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}

	user, token, err := h.Service.Login(payload)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			common.ErrorHandler(c, common.NewAppError(
				http.StatusUnauthorized,
				common.AUTH_INVALID_CREDENTIALS,
				err.Error(),
				nil,
			))
			return
		}
		var deletedErr *service.DeletedAccountError
		if errors.As(err, &deletedErr) {
			common.ErrorHandler(c, common.NewAppError(
				http.StatusForbidden,
				common.FORBIDDEN,
				err.Error(),
				map[string]string{
					"email":      deletedErr.Email,
					"deleted_at": deletedErr.DeletedAt.Format(time.RFC3339),
				},
			))
			return
		}
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		User: response.UserResponse{
			ID:        user.ID,
			FullName:  user.FullName,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, common.SuccessResponse(gin.H{
		"message": "logout success, remove bearer token on client",
	}))
}
