package common

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context, err error) {
	if err == nil {
		return
	}

	if appErr, ok := err.(*AppError); ok {
		log.Printf("request failed: method=%s path=%s code=%s status=%d err=%s", c.Request.Method, c.FullPath(), appErr.Code, appErr.StatusCode, appErr.Message)
		c.Errors = nil
		c.AbortWithStatusJSON(appErr.StatusCode, ErrorResponse[any](
			ErrorSchema{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			},
		))
		return
	}

	log.Printf("unhandled error: method=%s path=%s err=%v", c.Request.Method, c.FullPath(), err)
	serverErr := InternalServerError()
	c.Errors = nil
	c.AbortWithStatusJSON(serverErr.StatusCode, ErrorResponse[any](
		ErrorSchema{
			Code:    serverErr.Code,
			Message: serverErr.Message,
			Details: serverErr.Details,
		},
	))
}
