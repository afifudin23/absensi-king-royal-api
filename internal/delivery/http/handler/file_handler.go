package handler

import (
	"net/http"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	service     service.FileService
	userService service.UserService
}

func NewFileHandler(fileService service.FileService, userService service.UserService) *FileHandler {
	return &FileHandler{service: fileService, userService: userService}
}

func (h *FileHandler) Upload(c *gin.Context) {
	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		common.ErrorHandler(c, common.UnauthorizedError("User not authenticated"))
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		common.ErrorHandler(c, common.BadRequestError("File is required"))
		return
	}

	const maxSize = 5 << 20 // 5 MB
	if fileHeader.Size > maxSize {
		common.ErrorHandler(c, common.BadRequestError("File size exceeds 5 MB limit"))
		return
	}

	fileType := c.PostForm("file_type")
	if fileType == "" {
		common.ErrorHandler(c, common.BadRequestError("File type is required, must be one of: check_in, check_out, profile_picture, sick, extra_off, overtime, leave"))
		return
	}
	if !isValidFileType(fileType) {
		common.ErrorHandler(c, common.BadRequestError("Invalid file type: '"+fileType+"'. Allowed values: check_in, check_out, profile_picture, sick, extra_off, overtime, leave"))
		return
	}

	if model.FileType(fileType) == model.FileTypeProfilePicture && h.userService != nil {
		if user, err := h.userService.GetByID(c.Request.Context(), userID); err == nil {
			if user.Profile != nil && user.Profile.ProfilePictureID != nil {
				_ = h.service.Delete(c.Request.Context(), *user.Profile.ProfilePictureID)
			}
		}
	}

	file, err := h.service.Upload(c.Request.Context(), fileHeader, model.FileType(fileType), userID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	if model.FileType(fileType) == model.FileTypeProfilePicture && h.userService != nil {
		fileID := file.ID
		_, _ = h.userService.UpdateProfile(c.Request.Context(), userID, request.UserUpdateProfileRequest{
			ProfilePictureID: &fileID,
		})
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToFileResponse(*file)))
}

func (h *FileHandler) Delete(c *gin.Context) {
	fileID := c.Param("file_id")
	if fileID == "" {
		common.ErrorHandler(c, common.BadRequestError("File id is required"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), fileID); err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(fileID)))
}

func isValidFileType(t string) bool {
	switch model.FileType(t) {
	case model.FileTypeCheckIn, model.FileTypeCheckOut, model.FileTypeProfilePicture,
		model.FileTypeSick, model.FileTypeExtraOff, model.FileTypeOvertime, model.FileTypeLeave:
		return true
	default:
		return false
	}
}
