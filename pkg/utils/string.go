package utils

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/gin-gonic/gin"
)

func GetCurrentUserID(c *gin.Context) (string, bool) {
	uid, exists := c.Get("uid")
	if !exists {
		c.Error(common.UnauthorizedError("Unauthorized, please login again"))
		c.Abort()
		return "", false
	}

	userID, ok := uid.(string)
	if !ok || userID == "" {
		c.Error(common.UnauthorizedError("Unauthorized, please login again"))
		c.Abort()
		return "", false
	}

	return userID, true
}
