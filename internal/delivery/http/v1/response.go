package v1

import (
	"github.com/go-playground/validator/v10"
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

type dataValidation struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

type failureInfoValidation struct {
	Code int              `json:"code" example:"400"`
	Data []dataValidation `json:"error"  example:"invalid request body"`
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

func errorValidationResponse(c *gin.Context, err error) {

	data := []dataValidation{}

	for _, err := range err.(validator.ValidationErrors) {
		data = append(data, dataValidation{
			Key:     err.Field(),
			Message: err.Tag(),
		})
	}

	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, failureValidation{Error: failureInfoValidation{
		Code: http.StatusUnprocessableEntity,
		Data: data,
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
