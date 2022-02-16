package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type success struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
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
	context.JSON(http.StatusOK, success{Data: data, Message: "Success"})
}

func createdResponse(context *gin.Context, data interface{}) {
	context.JSON(http.StatusCreated, success{Data: data, Message: "Success"})
}

func errorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, failure{Error: failureInfo{
		Code:    statusCode,
		Message: message,
	}})
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
