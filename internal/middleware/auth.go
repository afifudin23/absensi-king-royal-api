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
			c.Error(common.UnauthorizedError("Authorization header is required"))
			c.Abort()
			return
		}

		// Require Bearer token format.
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.Error(common.UnauthorizedError("Authorization header must use Bearer token"))
			c.Abort()
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
		if token == "" {
			c.Error(common.UnauthorizedError("Bearer token is required"))
			c.Abort()
			return
		}

		// Verify token
		env := config.GetEnv()
		if env == nil {
			c.Error(common.InternalServerError())
			c.Abort()
			return
		}
		claims, err := utils.VerifyToken(token, env.AccessKey)
		if err != nil {
			c.Error(common.UnauthorizedError("Invalid or expired token"))
			c.Abort()
			return
		}

		// Check user still exists and is not soft-deleted.
		db := config.GetDB()
		if db == nil {
			c.Error(common.InternalServerError())
			c.Abort()
			return
		}

		var user model.User
		if err := db.Where("id = ? AND deleted_at IS NULL", claims.UID).Take(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Error(common.UnauthorizedError("User not found or has been deleted"))
				c.Abort()
				return
			}
			c.Error(common.InternalServerError())
			c.Abort()
			return
		}

		// Set user id from token claims for downstream handlers.
		c.Set("uid", user.ID)
		ctx := logger.WithUserID(c.Request.Context(), user.ID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
