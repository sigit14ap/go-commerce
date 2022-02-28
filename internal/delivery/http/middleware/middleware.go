package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/internal/service"
)

type Function interface {
	Handle(context *gin.Context)
}

type MiddlewareService struct {
	VerifyStore Function
}

func NewMiddleware(services *service.Services) *MiddlewareService {
	return &MiddlewareService{
		VerifyStore: NewVerifyStore(services),
	}
}
