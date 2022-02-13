package http

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/sigit14ap/go-commerce/internal/delivery/http/v1"
	"github.com/sigit14ap/go-commerce/internal/service"
	"github.com/sigit14ap/go-commerce/pkg/auth"
	"github.com/sigit14ap/go-commerce/pkg/storage"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"net/http"
	//_ "github.com/sigit14ap/go-commerce/docs"
)

type Handler struct {
	services        *service.Services
	tokenProvider   auth.TokenProvider
	storageProvider storage.StorageProvider
}

func NewHandler(services *service.Services, tokenProvider auth.TokenProvider, storageProvider storage.StorageProvider) *Handler {
	return &Handler{
		services:        services,
		tokenProvider:   tokenProvider,
		storageProvider: storageProvider,
	}
}

func (h *Handler) Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.tokenProvider, h.storageProvider)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
