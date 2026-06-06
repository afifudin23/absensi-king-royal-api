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

// Register godoc
//
//	@Summary		Registrasi pengguna baru
//	@Description	Mendaftarkan akun pengguna baru menggunakan email dan password.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `full_name` (string, required, min 3 karakter): Nama lengkap pengguna
//	@Description	- `email` (string, required, format email): Alamat email unik
//	@Description	- `password` (string, required, min 8 karakter): Password akun
//	@Description
//	@Description	**Validasi & Behavior:**
//	@Description	- Email harus unik, jika sudah terdaftar akan mengembalikan 400
//	@Description	- Password di-hash menggunakan argon2id sebelum disimpan
//	@Description
//	@Description	**Response 201:** ID pengguna yang baru dibuat
//	@Description
//	@Description	**Error:** 400 jika email sudah terdaftar atau validasi gagal
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.AuthRegisterRequest	true	"Data registrasi"
//	@Success		201		{object}	common.Response[response.RegisterResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		500		{object}	common.Response[any]
//	@Router			/auth/register [post]
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

// Login godoc
//
//	@Summary		Login pengguna
//	@Description	Login menggunakan email dan password, mengembalikan JWT access token.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `email` (string, required): Alamat email akun
//	@Description	- `password` (string, required): Password akun
//	@Description
//	@Description	**Validasi & Behavior:**
//	@Description	- Jika email tidak ditemukan atau password salah, request ditolak dengan 401
//	@Description	- Token yang dikembalikan digunakan sebagai Bearer token di header Authorization
//	@Description
//	@Description	**Response 200:** `access_token` (JWT), `token_type` (Bearer), data user (id, nama, email, role)
//	@Description
//	@Description	**Error:** 401 jika kredensial salah
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.AuthLoginRequest	true	"Kredensial login"
//	@Success		200		{object}	common.Response[response.LoginResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Router			/auth/login [post]
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

// Logout godoc
//
//	@Summary		Logout pengguna
//	@Description	Logout dari sesi aktif.
//	@Description
//	@Description	**Behavior:**
//	@Description	- Endpoint ini tidak menginvalidasi token di server (stateless JWT)
//	@Description	- Klien wajib menghapus `access_token` dari penyimpanan lokal setelah memanggil endpoint ini
//	@Description
//	@Description	**Response 200:** Pesan konfirmasi logout berhasil
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	common.Response[any]
//	@Router			/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, common.SuccessResponse(gin.H{
		"message": "logout success, remove bearer token on client",
	}))
}

// ForgotPassword godoc
//
//	@Summary		Lupa password
//	@Description	Meminta kode OTP yang dikirim ke email untuk proses reset password.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `email` (string, required, format email): Alamat email akun yang terdaftar
//	@Description
//	@Description	**Behavior:**
//	@Description	- Jika email terdaftar, kode OTP (6 digit) akan dikirim ke alamat email tersebut
//	@Description	- Jika email tidak terdaftar, response tetap 200 (untuk keamanan, tidak mengekspos data)
//	@Description	- Kode OTP memiliki masa berlaku terbatas
//	@Description
//	@Description	**Response 200:** Pesan konfirmasi pengiriman OTP
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.ForgotPasswordRequest	true	"Email terdaftar"
//	@Success		200		{object}	common.Response[any]
//	@Failure		400		{object}	common.Response[any]
//	@Router			/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var payload request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	if err := h.Service.ForgotPassword(c.Request.Context(), payload.Email); err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(gin.H{
		"message": "Jika email terdaftar, kode OTP akan dikirim ke email kamu.",
	}))
}

// ResetPassword godoc
//
//	@Summary		Reset password
//	@Description	Mereset password menggunakan kode OTP yang telah dikirim ke email.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `email` (string, required, format email): Email akun yang akan direset
//	@Description	- `otp` (string, required, tepat 6 karakter): Kode OTP dari email
//	@Description	- `new_password` (string, required, min 8 karakter): Password baru
//	@Description
//	@Description	**Validasi & Behavior:**
//	@Description	- OTP harus valid dan belum kadaluarsa
//	@Description	- Password baru akan di-hash sebelum disimpan
//	@Description	- Setelah berhasil, OTP akan dihapus dan tidak bisa digunakan lagi
//	@Description
//	@Description	**Response 200:** Pesan konfirmasi password berhasil direset
//	@Description
//	@Description	**Error:** 400 jika OTP tidak valid atau sudah kadaluarsa
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.ResetPasswordRequest	true	"Data reset password"
//	@Success		200		{object}	common.Response[any]
//	@Failure		400		{object}	common.Response[any]
//	@Router			/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var payload request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	if err := h.Service.ResetPassword(c.Request.Context(), payload); err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(gin.H{
		"message": "Password berhasil direset.",
	}))
}
