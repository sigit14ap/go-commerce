package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type success struct {
	Data interface{} `json:"data"`
}

type failure struct {
	Error failureInfo `json:"error"`
}

type failureValidation struct {
	Error failureInfoValidation `json:"error"`
}

type failureInfo struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"invalid request body"`
}

type failureInfoValidation struct {
	Code    int   `json:"code" example:"400"`
	Message error `json:"error"  example:"invalid request body"`
}

func successResponse(context *gin.Context, data interface{}) {
	log.Infof("Response with status OK: %v", data)
	context.JSON(http.StatusOK, success{Data: data})
}

func createdResponse(context *gin.Context, data interface{}) {
	log.Infof("Response with status Created: %v", data)
	context.JSON(http.StatusCreated, success{Data: data})
}

func errorResponse(c *gin.Context, statusCode int, message string) {
	// log.Error(message)
	c.AbortWithStatusJSON(statusCode, failure{Error: failureInfo{
		Code:    statusCode,
		Message: message,
	}})
}

func errorValidationResponse(c *gin.Context, statusCode int, err error) {
	// log.Error(message)
	c.AbortWithStatusJSON(statusCode, failureValidation{Error: failureInfoValidation{
		Code:    statusCode,
		Message: err,
	}})
}
