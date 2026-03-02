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

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{Service: service.NewAuthService()}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var payload request.AuthRegisterRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	user, err := h.Service.Register(payload)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyRegistered) {
			common.ErrorHandler(c, common.BadRequestError(err.Error()))
			return
		}
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response.ToRegisterResponse(user.ID)))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var payload request.AuthLoginRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

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

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToLoginResponse(response.UserData{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, token)))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, common.SuccessResponse(gin.H{
		"message": "logout success, remove bearer token on client",
	}))
}
