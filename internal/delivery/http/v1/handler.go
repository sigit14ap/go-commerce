package v1

import (
	"fmt"
	"github.com/sigit14ap/go-commerce/internal/delivery/http/middleware"
	"github.com/sigit14ap/go-commerce/internal/delivery/http/services"
	"github.com/sigit14ap/go-commerce/pkg/courier"
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
	courierProvider courier.CourierProvider
	middlewares     *middleware.MiddlewareService
}

func NewHandler(services *service.Services, tokenProvider auth.TokenProvider, storageProvider storage.StorageProvider, courierProvider courier.CourierProvider, middlewares *middleware.MiddlewareService) *Handler {
	return &Handler{
		services:        services,
		tokenProvider:   tokenProvider,
		storageProvider: storageProvider,
		courierProvider: courierProvider,
		middlewares:     middlewares,
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

		user := v1.Group("user")
		{
			h.initUserAuthRoutes(user)

			v1.Use(h.verifyUser)
			{
				h.initUserAddressRoutes(v1)
			}
		}

		v1.Use(h.verifyUser)
		{
			store := v1.Group("store")
			{
				h.initStoreRoutes(store)

				v1.Use(h.middlewares.VerifyStore.Handle)
				{
					h.initStoreProductRoutes(store)
				}
			}
		}
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
	return services.GetIdFromPath(c, paramName)
}

func getIdFromRequestContext(context *gin.Context, paramName string) (primitive.ObjectID, error) {
	return services.GetIdFromRequestContext(context, paramName)
}

func getIdFromRequest(paramName string) (primitive.ObjectID, error) {
	return services.GetIdFromRequest(paramName)
}
