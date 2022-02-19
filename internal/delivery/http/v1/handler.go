package v1

import (
	"errors"
	"fmt"
	"github.com/sigit14ap/go-commerce/pkg/storage"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/internal/service"
	"github.com/sigit14ap/go-commerce/pkg/auth"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	v1.Use(LoggerMiddleware())
	{
		h.initAdminsRoutes(v1)
		h.initUsersRoutes(v1)
		h.initProductsRoutes(v1)
		h.initCartRoutes(v1)
		h.initOrdersRoutes(v1)
		h.initAreasRoutes(v1)
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()

		if len(c.Errors) > 0 {
			log.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("[%s %d] %s (%dms)", c.Request.Method, statusCode, path, latency)
			if statusCode >= http.StatusInternalServerError {
				log.Error(msg)
			} else if statusCode >= http.StatusBadRequest {
				log.Warn(msg)
			} else {
				log.Info(msg)
			}
		}
	}
}

func getIdFromPath(c *gin.Context, paramName string) (primitive.ObjectID, error) {
	idString := c.Param(paramName)
	if idString == "" {
		return primitive.ObjectID{}, errors.New("empty id param")
	}

	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return primitive.ObjectID{}, errors.New("invalid id param")
	}

	return id, nil
}

func getIdFromRequestContext(context *gin.Context, paramName string) (primitive.ObjectID, error) {
	idString, ok := context.Get(paramName)
	if !ok {
		return primitive.ObjectID{}, errors.New("not authenticated")
	}

	id, err := primitive.ObjectIDFromHex(idString.(string))
	if err != nil {
		return primitive.ObjectID{}, errors.New("invalid id param")
	}

	return id, nil
}

func getIdFromRequest(paramName string) (primitive.ObjectID, error) {
	if paramName == "" {
		return primitive.ObjectID{}, errors.New("empty id param")
	}

	id, err := primitive.ObjectIDFromHex(paramName)
	if err != nil {
		return primitive.ObjectID{}, errors.New("invalid id param")
	}

	return id, nil
}
