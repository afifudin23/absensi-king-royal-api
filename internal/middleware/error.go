package middleware

import (
	"log"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("PANIC:", r)

				common.ErrorHandler(c, common.InternalServerError())

				c.Abort()
			}
		}()

		c.Next()
	}
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 && !c.Writer.Written() {
			common.ErrorHandler(c, c.Errors.Last().Err)
			c.Abort()
		}
	}
}
