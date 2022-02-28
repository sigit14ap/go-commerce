package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/internal/delivery/http/services"
)

func successResponse(context *gin.Context, data interface{}) {
	services.SuccessResponse(context, data)
}

func createdResponse(context *gin.Context, data interface{}) {
	services.CreatedResponse(context, data)
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	services.ErrorResponse(c, statusCode, message)
}

func errorValidationResponse(c *gin.Context, err error) {
	services.ErrorValidationResponse(c, err)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
