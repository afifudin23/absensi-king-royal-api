package handler

import (
	"errors"
	"net/http"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{Service: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var payload request.AuthRegisterRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	user, err := h.Service.Register(c.Request.Context(), payload)
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

	user, token, err := h.Service.Login(c.Request.Context(), payload)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			common.ErrorHandler(c, common.UnauthorizedError(err.Error()))
			return
		}
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToLoginResponse(*user, token)))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, common.SuccessResponse(gin.H{
		"message": "logout success, remove bearer token on client",
	}))
}
