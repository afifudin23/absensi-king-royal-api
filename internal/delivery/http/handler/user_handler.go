package handler

import (
	"net/http"
	"strings"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
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

// GetAllUsers godoc
//
//	@Summary		Daftar semua pengguna
//	@Description	Mendapatkan daftar seluruh pengguna yang terdaftar.
//	@Description
//	@Description	**Query Params (opsional):**
//	@Description	- `search` (string): Filter berdasarkan nama atau email
//	@Description	- `role` (string): Filter berdasarkan role — `admin` atau `user`
//	@Description
//	@Description	**Response 200:** Array data pengguna beserta profil lengkap (posisi, departemen, gaji, dll)
//	@Tags			Users
//	@Produce		json
//	@Param			search	query		string	false	"Cari berdasarkan nama atau email"
//	@Param			role	query		string	false	"Filter berdasarkan role (admin/user)"
//	@Success		200		{object}	common.Response[[]response.UserResponse]
//	@Failure		500		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	filter := &repository.UserFilter{
		Search: strings.TrimSpace(c.Query("search")),
		Role:   strings.TrimSpace(c.Query("role")),
	}
	users, err := h.Service.GetAll(c.Request.Context(), filter)
	if err != nil {
		common.ErrorHandler(c, common.InternalServerError())
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToUserListResponse(users)))
}

// ResetUserPassword godoc
//
//	@Summary		Reset password pengguna (Admin)
//	@Description	Admin mereset password pengguna tertentu tanpa perlu mengetahui password lama.
//	@Description
//	@Description	**Path Param:**
//	@Description	- `user_id` (UUID): ID pengguna yang akan direset passwordnya
//	@Description
//	@Description	**Request Body:**
//	@Description	- `new_password` (string, required, min 8 karakter): Password baru untuk pengguna tersebut
//	@Description
//	@Description	**Akses:** Hanya admin
//	@Description
//	@Description	**Response 200:** ID pengguna yang passwordnya berhasil direset
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		string								true	"ID pengguna"
//	@Param			body	body		request.AdminResetPasswordRequest	true	"Password baru"
//	@Success		200		{object}	common.Response[common.ActionSuccessResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Failure		403		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/users/{user_id}/reset-password [post]
func (h *UserHandler) ResetUserPassword(c *gin.Context) {
	var payload request.AdminResetPasswordRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}

	adminID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}
	targetUserID := c.Param("user_id")

	if err := h.Service.ResetPassword(c.Request.Context(), adminID, targetUserID, payload.NewPassword); err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(targetUserID)))
}

// GetMyProfile godoc
//
//	@Summary		Profil saya
//	@Description	Mendapatkan data profil lengkap pengguna yang sedang login.
//	@Description
//	@Description	**Response 200:** Data pengguna lengkap termasuk profil karyawan (posisi, departemen, gaji pokok, foto profil, dll)
//	@Description
//	@Description	**Error:** 401 jika token tidak valid atau tidak ada
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	common.Response[response.UserResponse]
//	@Failure		401	{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/users/me [get]
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

// UpdateMyProfile godoc
//
//	@Summary		Update profil saya
//	@Description	Memperbarui data profil pengguna yang sedang login. Semua field bersifat opsional (partial update).
//	@Description
//	@Description	**Field yang bisa diupdate:**
//	@Description	- `full_name`, `email`, `phone_number`, `address`, `gender`, `birth_date`, `birth_place`
//	@Description	- `employee_code`, `employment_status`, `position`, `department`, `bank_account_number`
//	@Description	- `profile_picture_id` (UUID file yang sudah diupload), `clear_profile_picture` (boolean, untuk hapus foto)
//	@Description
//	@Description	**Response 200:** ID pengguna yang diupdate
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.UserUpdateProfileRequest	true	"Data profil yang diperbarui"
//	@Success		200		{object}	common.Response[common.ActionSuccessResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/users/me [put]
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

// CreateUser godoc
//
//	@Summary		Buat pengguna baru
//	@Description	Membuat akun pengguna baru beserta data profil karyawan.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `full_name` (string, required): Nama lengkap
//	@Description	- `email` (string, required): Email unik
//	@Description	- `role` (string, required): Role akun — `admin` atau `user`
//	@Description	- `employee_code`, `employment_status`, `birth_place`, `birth_date`, `gender` (opsional)
//	@Description	- `address`, `phone_number`, `position`, `department`, `bank_account_number` (opsional)
//	@Description	- `basic_salary`, `position_allowance`, `other_allowance` (opsional, angka)
//	@Description
//	@Description	**Behavior:**
//	@Description	- Password awal akan di-generate otomatis dan bisa direset oleh admin
//	@Description	- Mengembalikan daftar semua pengguna setelah user baru berhasil dibuat
//	@Description
//	@Description	**Response 201:** Array seluruh data pengguna (termasuk yang baru dibuat)
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.UserCreateRequest	true	"Data pengguna baru"
//	@Success		201		{object}	common.Response[[]response.UserResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var payload request.UserCreateRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	_, err := h.Service.Create(c.Request.Context(), payload)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	users, err := h.Service.GetAll(c.Request.Context(), nil)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusCreated, common.SuccessResponse(response.ToUserListResponse(users)))
}

// GetUserByID godoc
//
//	@Summary		Detail pengguna
//	@Description	Mendapatkan data lengkap satu pengguna berdasarkan ID.
//	@Description
//	@Description	**Path Param:**
//	@Description	- `user_id` (UUID): ID pengguna yang ingin dilihat
//	@Description
//	@Description	**Response 200:** Data pengguna beserta profil lengkap
//	@Description
//	@Description	**Error:** 404 jika pengguna tidak ditemukan
//	@Tags			Users
//	@Produce		json
//	@Param			user_id	path		string	true	"ID pengguna"
//	@Success		200		{object}	common.Response[response.UserResponse]
//	@Failure		401		{object}	common.Response[any]
//	@Failure		404		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/users/{user_id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("user_id")
	user, err := h.Service.GetByID(c.Request.Context(), userID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToUserResponse(*user)))
}

// UpdateUser godoc
//
//	@Summary		Update pengguna
//	@Description	Memperbarui data pengguna tertentu berdasarkan ID. Semua field opsional (partial update).
//	@Description
//	@Description	**Field yang bisa diupdate:** `full_name`, `email`, `role`, `employee_code`, `employment_status`,
//	@Description	`birth_place`, `birth_date`, `gender`, `address`, `phone_number`, `position`, `department`,
//	@Description	`bank_account_number`, `basic_salary`, `position_allowance`, `other_allowance`, `joined_at`,
//	@Description	`profile_picture_id`, `clear_profile_picture`
//	@Description
//	@Description	**Response 200:** Array seluruh data pengguna setelah diupdate
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		string						true	"ID pengguna"
//	@Param			body	body		request.UserUpdateRequest	true	"Data yang diperbarui"
//	@Success		200		{object}	common.Response[[]response.UserResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/users/{user_id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var payload request.UserUpdateRequest
	userID := c.Param("user_id")

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	_, err := h.Service.Update(c.Request.Context(), userID, payload)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	users, err := h.Service.GetAll(c.Request.Context(), nil)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToUserListResponse(users)))
}

// DeleteUser godoc
//
//	@Summary		Hapus pengguna
//	@Description	Menghapus pengguna berdasarkan ID menggunakan soft delete (data tidak benar-benar terhapus dari database).
//	@Description
//	@Description	**Path Param:**
//	@Description	- `user_id` (UUID): ID pengguna yang akan dihapus
//	@Description
//	@Description	**Response 200:** ID pengguna yang dihapus
//	@Description
//	@Description	**Error:** 404 jika pengguna tidak ditemukan
//	@Tags			Users
//	@Produce		json
//	@Param			user_id	path		string	true	"ID pengguna"
//	@Success		200		{object}	common.Response[common.ActionSuccessResponse]
//	@Failure		401		{object}	common.Response[any]
//	@Failure		404		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/users/{user_id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	err := h.Service.Delete(c.Request.Context(), userID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(userID)))
}

// ChangePassword godoc
//
//	@Summary		Ganti password
//	@Description	Mengganti password akun pengguna yang sedang login.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `current_password` (string, required): Password saat ini
//	@Description	- `new_password` (string, required, min 8 karakter): Password baru
//	@Description	- `confirm_password` (string, required): Konfirmasi password baru, harus sama dengan `new_password`
//	@Description
//	@Description	**Validasi:**
//	@Description	- `current_password` harus cocok dengan password yang tersimpan
//	@Description	- `new_password` dan `confirm_password` harus identik
//	@Description
//	@Description	**Response 200:** ID pengguna
//	@Description
//	@Description	**Error:** 400 jika password lama salah atau konfirmasi tidak cocok
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.ChangePasswordRequest	true	"Password lama dan baru"
//	@Success		200		{object}	common.Response[common.ActionSuccessResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/users/me/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var payload request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}

	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}

	if err := h.Service.ChangePassword(c.Request.Context(), userID, payload); err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(userID)))
}
