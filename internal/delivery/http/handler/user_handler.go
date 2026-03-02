package handler

import (
	"net/http"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{Service: service.NewUserService()}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		common.ErrorHandler(c, common.InternalServerError())
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToUserListResponse(users)))
}

func (h *UserHandler) GetMyProfile(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		c.Error(common.UnauthorizedError("Unauthorized, please login again"))
		c.Abort()
		return
	}
	user, err := h.Service.GetUserByID(uid.(string))
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

	user, err := h.Service.UpdateUser(uid.(string), model.User{
		FullName:          stringValue(payload.FullName),
		Email:             stringValue(payload.Email),
		Password:          stringValue(payload.Password),
		Role:              stringValue(payload.Role),
		EmployeeCode:      payload.EmployeeCode,
		EmploymentStatus:  payload.EmploymentStatus,
		BirthPlace:        payload.BirthPlace,
		BirthDate:         payload.BirthDate,
		Gender:            payload.Gender,
		Address:           payload.Address,
		PhoneNumber:       payload.PhoneNumber,
		Position:          payload.Position,
		Department:        payload.Department,
		BankAccountNumber: payload.BankAccountNumber,
		ProfilePictureURL: payload.ProfilePictureURL,
		ProfilePictureID:  payload.ProfilePictureID,
	})

	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToUserSuccessResponse(user.ID)))
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var payload request.UserCreateRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	user, err := h.Service.CreateUser(model.User{
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: payload.Password,
		Role:     payload.Role,

		EmployeeCode:      payload.EmployeeCode,
		EmploymentStatus:  payload.EmploymentStatus,
		BirthPlace:        payload.BirthPlace,
		BirthDate:         payload.BirthDate,
		Gender:            payload.Gender,
		Address:           payload.Address,
		PhoneNumber:       payload.PhoneNumber,
		Position:          payload.Position,
		Department:        payload.Department,
		BankAccountNumber: payload.BankAccountNumber,
		ProfilePictureURL: payload.ProfilePictureURL,
		ProfilePictureID:  payload.ProfilePictureID,
	})
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusCreated, common.SuccessResponse(response.ToUserSuccessResponse(user.ID)))
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("user_id")
	user, err := h.Service.GetUserByID(userID)
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

	user, err := h.Service.UpdateUser(userID, model.User{
		FullName:          stringValue(payload.FullName),
		Role:              stringValue(payload.Role),
		EmployeeCode:      payload.EmployeeCode,
		EmploymentStatus:  payload.EmploymentStatus,
		BirthPlace:        payload.BirthPlace,
		BirthDate:         payload.BirthDate,
		Gender:            payload.Gender,
		Address:           payload.Address,
		PhoneNumber:       payload.PhoneNumber,
		Position:          payload.Position,
		Department:        payload.Department,
		BankAccountNumber: payload.BankAccountNumber,
		ProfilePictureURL: payload.ProfilePictureURL,
		ProfilePictureID:  payload.ProfilePictureID,
	})

	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToUserSuccessResponse(user.ID)))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	err := h.Service.DeleteUser(userID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToUserSuccessResponse(userID)))
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
