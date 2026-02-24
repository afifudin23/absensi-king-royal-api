package common

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Custom application errors
	if appErr, ok := err.(*AppError); ok {
		log.Printf("request failed: method=%s path=%s code=%s status=%d err=%s", c.Request.Method, c.FullPath(), appErr.Code, appErr.StatusCode, appErr.Message)
		c.JSON(appErr.StatusCode, ErrorResponse[any](
			ErrorSchema{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			},
		))
		c.Errors = nil
		return
	}

	// Fallback internal server error response
	log.Printf("unhandled error: method=%s path=%s err=%v", c.Request.Method, c.FullPath(), err)
	serverErr := InternalServerError()
	c.JSON(serverErr.StatusCode, ErrorResponse[any](
		ErrorSchema{
			Code:    serverErr.Code,
			Message: serverErr.Message,
			Details: serverErr.Details,
		},
	))
	c.Errors = nil
}
