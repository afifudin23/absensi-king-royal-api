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
	log.Println(err.Error())
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
