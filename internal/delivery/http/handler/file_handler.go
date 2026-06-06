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

// Upload godoc
//
//	@Summary		Upload file
//	@Description	Upload file ke server untuk digunakan pada check-in, check-out, foto profil, atau bukti pengajuan.
//	@Description
//	@Description	**Form Data:**
//	@Description	- `file` (file, required): File yang akan diupload. Maksimal 5 MB
//	@Description	- `file_type` (string, required): Tipe file — `check_in`, `check_out`, `profile_picture`, `sick`, `extra_off`, `overtime`, `leave`
//	@Description
//	@Description	**Validasi:**
//	@Description	- Ukuran file maksimal 5 MB
//	@Description	- Tipe file harus salah satu dari nilai yang diizinkan di atas
//	@Description	- Jika tipe `profile_picture`, foto profil lama akan otomatis dihapus dan profil pengguna diperbarui
//	@Description
//	@Description	**Response 200:** Data file yang berhasil diupload (id, url, nama, ukuran, mime type)
//	@Description
//	@Description	**Cara pakai:** Upload file terlebih dahulu, lalu gunakan `id` yang dikembalikan di endpoint lain (check-in, pengajuan, dll)
//	@Tags			Files
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file		formData	file	true	"File yang akan diupload (maks 5 MB)"
//	@Param			file_type	formData	string	true	"Tipe file: check_in | check_out | profile_picture | sick | extra_off | overtime | leave"
//	@Success		200			{object}	common.Response[response.FileResponse]
//	@Failure		400			{object}	common.Response[any]
//	@Failure		401			{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/files [post]
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

// Delete godoc
//
//	@Summary		Hapus file
//	@Description	Menghapus file dari server berdasarkan ID.
//	@Description
//	@Description	**Path Param:**
//	@Description	- `file_id` (UUID): ID file yang akan dihapus
//	@Description
//	@Description	**Behavior:** File akan dihapus permanen dari server dan database
//	@Description
//	@Description	**Response 200:** ID file yang berhasil dihapus
//	@Tags			Files
//	@Produce		json
//	@Param			file_id	path		string	true	"ID file"
//	@Success		200		{object}	common.Response[common.ActionSuccessResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/files/{file_id} [delete]
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
