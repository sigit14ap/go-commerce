package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
	"github.com/sigit14ap/go-commerce/internal/delivery/http/services"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"path/filepath"
)

//TODO: add product search by query
func (h *Handler) initStoreProductRoutes(api *gin.RouterGroup) {
	products := api.Group("/products")
	{
		products.GET("/", h.storeGetProduct)
		products.GET("/:id", h.storeDetailProduct)
		products.POST("/", h.storeCreateProduct)
		products.PUT("/:id", h.storeUpdateProduct)
		products.DELETE("/:id", h.storeDeleteProduct)
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

// StoreCreateProduct godoc
// @Summary   Create product store
// @Tags      store-products
// @Accept    json
// @Produce   json
// @Param     product  body      dto.CreateProductDTO  true  "product"
// @Success   201      {object}  success
// @Failure   400      {object}  failure
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  StoreAuth
// @Router    /store/products [post]
func (h *Handler) storeCreateProduct(context *gin.Context) {

	storeID, _ := services.GetIdFromRequestContext(context, "storeID")

	var productInput dto.CreateProductInput
	err := context.ShouldBindWith(&productInput, binding.FormMultipart)
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, "invalid input body")
		return
	}

	categoryID, err := primitive.ObjectIDFromHex(productInput.CategoryID)
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	category, err := h.services.Categories.FindByID(context.Request.Context(), categoryID)
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	form, _ := context.MultipartForm()
	files := form.File["images[][image]"]

	if len(files) == 0 {
		ErrorResponse(context, http.StatusBadRequest, "images required")
		return
	}

	var imagesArray []string

	for _, file := range files {

		allowedExt := []string{".jpg", ".png", "jpeg"}
		isContains := contains(allowedExt, filepath.Ext(file.Filename))

		if !isContains {
			ErrorResponse(context, http.StatusBadRequest, "Icon must be jpg, jpeg or png")
			return
		}

		uploadedFile := h.storageProvider.Upload("Category", file)
		imagesArray = append(imagesArray, uploadedFile)
	}

	productDTO := dto.CreateProductDTO{}
	copier.Copy(&productDTO, &productInput)
	productDTO.CategoryID = category.ID
	productDTO.StoreID = storeID

	productDTO.Images = imagesArray
	product, err := h.services.Products.Create(context.Request.Context(), productDTO)
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, product)
}

// StoreUpdateProduct godoc
// @Summary   Update product store
// @Tags      store-products
// @Accept    json
// @Produce   json
// @Param     id       path      string                true  "product id"
// @Param     product  body      dto.UpdateProductDTO  true  "product update fields"
// @Success   200  {object}  success
// @Failure   400  {object}  failure
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  StoreAuth
// @Router    /store/products/{id} [put]
func (h *Handler) storeUpdateProduct(context *gin.Context) {

	storeID, _ := services.GetIdFromRequestContext(context, "storeID")

	var productInput dto.UpdateProductInput

	err := context.ShouldBindWith(&productInput, binding.FormMultipart)
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, "invalid input body")
		return
	}

	productID, err := getIdFromPath(context, "id")
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	categoryID, err := primitive.ObjectIDFromHex(productInput.CategoryID)
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	category, err := h.services.Categories.FindByID(context.Request.Context(), categoryID)
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info(context.Request.PostForm["images"])
	form, _ := context.MultipartForm()
	files := form.File["images[]"]
	log.Info(files)
	//if len(files) == 0 {
	//	errorResponse(context, http.StatusBadRequest, "images required")
	//	return
	//}
	//
	//var imagesArray []string
	//
	//for _, file := range files {
	//
	//	allowedExt := []string{".jpg", ".png", "jpeg"}
	//	isContains := contains(allowedExt, filepath.Ext(file.Filename))
	//
	//	if !isContains {
	//		errorResponse(context, http.StatusBadRequest, "Icon must be jpg, jpeg or png")
	//		return
	//	}
	//
	//	uploadedFile := h.storageProvider.Upload("Category", file)
	//	imagesArray = append(imagesArray, uploadedFile)
	//}

	productDTO := dto.UpdateProductDTO{}
	copier.Copy(&productDTO, &productInput)
	productDTO.CategoryID = category.ID
	productDTO.StoreID = storeID

	//productDTO.Images = imagesArray

	product, err := h.services.Products.Update(context.Request.Context(), productDTO, productID)
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, product)
}

// StoreDeleteProduct godoc
// @Summary   Delete product store
// @Tags      store-products
// @Accept    json
// @Produce   json
// @Param     id   path      string  true  "product id"
// @Success   200  {object}  success
// @Failure   400  {object}  failure
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  StoreAuth
// @Router    /store/products/{id} [delete]
func (h *Handler) storeDeleteProduct(context *gin.Context) {

	storeID, _ := services.GetIdFromRequestContext(context, "storeID")

	productID, err := getIdFromPath(context, "id")
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Products.Delete(context, productID, storeID)
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	var data interface{}
	successResponse(context, data)
}
