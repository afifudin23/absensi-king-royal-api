package handler

import (
	"net/http"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{Service: userService}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.Service.GetAll(c.Request.Context())
	if err != nil {
		common.ErrorHandler(c, common.InternalServerError())
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToUserListResponse(users)))
}

func (h *UserHandler) GetMyProfile(c *gin.Context) {
	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}
	user, err := h.Service.GetByID(c.Request.Context(), userID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToUserResponse(*user)))
}

func (h *UserHandler) UpdateMyProfile(c *gin.Context) {
	var payload request.UserUpdateProfileRequest
	uid, exists := c.Get("uid")
	if !exists {
		c.Error(common.UnauthorizedError("Unauthorized, please login again"))
		c.Abort()
		return
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	user, err := h.Service.UpdateProfile(c.Request.Context(), uid.(string), payload)

	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(user.ID)))
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var payload request.UserCreateRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	user, err := h.Service.Create(c.Request.Context(), payload)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusCreated, common.SuccessResponse(common.ToSuccessResponse(user.ID)))
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("user_id")
	user, err := h.Service.GetByID(c.Request.Context(), userID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToUserResponse(*user)))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var payload request.UserUpdateRequest
	userID := c.Param("user_id")

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	user, err := h.Service.Update(c.Request.Context(), userID, payload)

	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(user.ID)))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	err := h.Service.Delete(c.Request.Context(), userID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(userID)))
}
