package middleware

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			common.ErrorHandler(c, common.UnauthorizedError("Unauthorized, please login again"))
			return
		}

		userRole, ok := role.(model.UserRole)
		if !ok || userRole == "" {
			common.ErrorHandler(c, common.UnauthorizedError("Unauthorized, please login again"))
			return
		}

		if userRole != model.UserRoleAdmin {
			common.ErrorHandler(c, common.ForbiddenError("Forbidden"))
			return
		}

		c.Next()
	}
}
