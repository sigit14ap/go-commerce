package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

//TODO: add product search by query
func (h *Handler) initStoreProductRoutes(api *gin.RouterGroup) {
	products := api.Group("/products")
	{
		products.GET("/", h.storeGetProduct)
		products.GET("/:id", h.storeDetailProduct)
		products.POST("/", h.createProductAdmin)
		products.PUT("/:id", h.updateProductAdmin)
		products.DELETE("/:id", h.deleteProductAdmin)
		products.GET("/:id/reviews", h.getProductReviewsAdmin)
	}
}

// StoreGetProduct godoc
// @Summary  Get all products store
// @Tags     store-products
// @Accept   json
// @Produce  json
// @Success  200  {array}   success
// @Failure  401  {object}  failure
// @Failure  404  {object}  failure
// @Failure  500  {object}  failure
// @Security  StoreAuth
// @Router   /store/products [get]
func (h *Handler) storeGetProduct(context *gin.Context) {
	products, err := h.services.Products.FindAll(context.Request.Context())
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	productsArray := make([]domain.Product, len(products))
	if products != nil {
		productsArray = products
	}

	successResponse(context, productsArray)
}

// StoreDetailProduct godoc
// @Summary   Get product by id
// @Tags      store-products
// @Accept    json
// @Produce   json
// @Param     id   path      string  true  "product id"
// @Success   200  {object}  success
// @Failure   400  {object}  failure
// @Failure   401      {object}  failure
// @Failure   404      {object}  failure
// @Failure   500      {object}  failure
// @Security  StoreAuth
// @Router    /store/products/{id} [get]
func (h *Handler) storeDetailProduct(context *gin.Context) {
	id, err := getIdFromPath(context, "id")
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}
	product, err := h.services.Products.FindByID(context.Request.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ErrorResponse(context, http.StatusInternalServerError,
				fmt.Sprintf("no products with id: %s", id.Hex()))
		} else {
			ErrorResponse(context, http.StatusInternalServerError, err.Error())
		}
		return
	}

	successResponse(context, product)
}
