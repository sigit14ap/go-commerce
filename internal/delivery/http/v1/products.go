package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

//TODO: add product search by query
func (h *Handler) initProductsRoutes(api *gin.RouterGroup) {
	products := api.Group("/products")
	{
		products.GET("/", h.getAllProducts)
		products.GET("/:id", h.getProductById)
		products.GET("/:id/reviews", h.getProductReviews)

		authenticated := products.Group("/", h.verifyUser)
		{
			authenticated.POST("/:id/reviews", h.createProductReview)
		}
	}
}

// GetProducts godoc
// @Summary  Get all products
// @Tags     products
// @Accept   json
// @Produce  json
// @Success  200  {array}   success
// @Failure  401  {object}  failure
// @Failure  404  {object}  failure
// @Failure  500  {object}  failure
// @Router   /products [get]
func (h *Handler) getAllProducts(context *gin.Context) {
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

// GetProductById godoc
// @Summary  Get product by id
// @Tags     products
// @Accept   json
// @Produce  json
// @Param    id   path      string  true  "product id"
// @Success  200  {object}  success
// @Failure  400  {object}  failure
// @Failure  401  {object}  failure
// @Failure  404  {object}  failure
// @Failure  500  {object}  failure
// @Router   /products/{id} [get]
func (h *Handler) getProductById(context *gin.Context) {
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

// GetProductReviews godoc
// @Summary  Get product reviews list
// @Tags     products
// @Accept   json
// @Produce  json
// @Param    id   path      string  true  "product id"
// @Success  200  {object}  success
// @Failure  400  {object}  failure
// @Failure  401  {object}  failure
// @Failure  404  {object}  failure
// @Failure  500  {object}  failure
// @Router   /products/{id}/reviews [get]
func (h *Handler) getProductReviews(context *gin.Context) {
	productID, err := getIdFromPath(context, "id")
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	reviews, err := h.services.Reviews.FindByProductID(context, productID)
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, reviews)
}

// CreateReview godoc
// @Summary   Create review
// @Tags      products
// @Accept    json
// @Produce   json
// @Param     id      path      string                   true  "product id"
// @Param     review  body      dto.CreateReviewDTOUser  true  "review"
// @Success   201     {object}  success
// @Failure   400     {object}  failure
// @Failure   401     {object}  failure
// @Failure   404     {object}  failure
// @Failure   500     {object}  failure
// @Security  UserAuth
// @Router    /products/{id}/reviews [post]
func (h *Handler) createProductReview(context *gin.Context) {
	productID, err := getIdFromPath(context, "id")
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := getIdFromRequestContext(context, "userID")
	if err != nil {
		ErrorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	var createDTO dto.CreateReviewDTOUser
	err = context.BindJSON(&createDTO)
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, "invalid input body")
		return
	}

	review, err := h.services.Reviews.Create(context, dto.CreateReviewInput{
		UserID:    userID,
		ProductID: productID,
		Text:      createDTO.Text,
		Rating:    createDTO.Rating,
	})

	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	createdResponse(context, review)
}

// GetProductsAdmin godoc
// @Summary   Get all products
// @Tags      admin-products
// @Accept    json
// @Produce   json
// @Success   200  {array}   success
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  AdminAuth
// @Router    /admins/products [get]
func (h *Handler) getAllProductsAdmin(context *gin.Context) {
	h.getAllProducts(context)
}

// GetProductByIdAdmin godoc
// @Summary   Get product by id
// @Tags      admin-products
// @Accept    json
// @Produce   json
// @Param     id   path      string  true  "product id"
// @Success   200  {object}  success
// @Failure   400  {object}  failure
// @Failure   401      {object}  failure
// @Failure   404      {object}  failure
// @Failure   500      {object}  failure
// @Security  AdminAuth
// @Router    /admins/products/{id} [get]
func (h *Handler) getProductByIdAdmin(context *gin.Context) {
	h.getProductById(context)
}

// GetProductReviewsAdmin godoc
// @Summary   Get product reviews list
// @Tags      admin-products
// @Accept    json
// @Produce   json
// @Param     id   path      string  true  "product id"
// @Success   200      {object}  success
// @Failure   400      {object}  failure
// @Failure   401      {object}  failure
// @Failure   404      {object}  failure
// @Failure   500      {object}  failure
// @Security  AdminAuth
// @Router    /admins/products/{id}/reviews [get]
func (h *Handler) getProductReviewsAdmin(context *gin.Context) {
	h.getProductReviews(context)
}
