package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/internal/delivery/http/services"
	"github.com/sigit14ap/go-commerce/internal/service"
	"net/http"
)

type VerifyStore struct {
	Handler *service.Services
}

func NewVerifyStore(services *service.Services) *VerifyStore {
	return &VerifyStore{
		Handler: services,
	}
}

func (verify *VerifyStore) Handle(context *gin.Context) {
	userID, err := services.GetIdFromRequestContext(context, "userID")

	if err != nil {
		services.ErrorResponse(context, http.StatusUnauthorized, "Unauthorized")
		return
	}

	store, err := verify.Handler.Stores.FindByUserID(context, userID)

	if err != nil {
		services.ErrorResponse(context, http.StatusForbidden, "Not registered as store")
		return
	}

	context.Set("storeID", store.ID.Hex())
	context.Set("storeData", store)
}
