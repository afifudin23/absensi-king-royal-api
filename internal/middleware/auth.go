package middleware

import (
	"errors"
	"strings"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/pkg/logger"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.ErrorHandler(c, common.UnauthorizedError("Authorization header is required"))
			return
		}

		// Require Bearer token format.
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			common.ErrorHandler(c, common.UnauthorizedError("Authorization header must use Bearer token"))
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
		if token == "" {
			common.ErrorHandler(c, common.UnauthorizedError("Bearer token is required"))
			return
		}

		// Verify token
		env := config.GetEnv()
		if env == nil {
			common.ErrorHandler(c, common.InternalServerError())
			return
		}
		claims, err := utils.VerifyToken(token, env.AccessKey)
		if err != nil {
			common.ErrorHandler(c, common.UnauthorizedError("Invalid or expired token"))
			return
		}

		// Check user still exists and is not soft-deleted.
		db := config.GetDB()
		if db == nil {
			common.ErrorHandler(c, common.InternalServerError())
			return
		}

		var user model.User
		if err := db.Where("id = ?", claims.UID).Take(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				common.ErrorHandler(c, common.UnauthorizedError("User not found or has been deleted"))
				return
			}
			common.ErrorHandler(c, common.InternalServerError())
			return
		}

		// Set user id from token claims for downstream handlers.
		c.Set("uid", user.ID)
		ctx := logger.WithUserID(c.Request.Context(), user.ID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
