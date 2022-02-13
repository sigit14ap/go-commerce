package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
)

var validate *validator.Validate = validator.New()

// GetProducts godoc
// @Summary  Get all category
// @Tags     admin-categories
// @Accept   json
// @Produce  json
// @Success  200  {array}   success
// @Failure  401  {object}  failure
// @Failure  404  {object}  failure
// @Failure  500  {object}  failure
// @Router   /admins/categories [get]
func (h *Handler) getAllCategoryAdmin(context *gin.Context) {
	categories, err := h.services.Categories.FindAll(context.Request.Context())
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	categoriesArray := make([]domain.Category, len(categories))
	if categories != nil {
		categoriesArray = categories
	}

	successResponse(context, categoriesArray)
}

// CreateProduct godoc
// @Summary   Create category
// @Tags      admin-category
// @Accept    json
// @Produce   json
// @Param     category  body      dto.CreateCategoryDTO  true  "category"
// @Success   201      {object}  success
// @Failure   400      {object}  failure
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  AdminAuth
// @Router    /admins/categories [post]
func (h *Handler) createCategoryAdmin(context *gin.Context) {
	var categoryDTO dto.CreateCategoryDTO
	err := context.ShouldBindWith(&categoryDTO, binding.FormMultipart)

	if err != nil {
		errorResponse(context, http.StatusUnprocessableEntity, "invalid input body")
		return
	}

	icon, err := context.FormFile("icon")
	if err != nil {
		log.Fatal(err)
	}

	allowedExt := []string{".jpg", ".png", "jpeg"}
	isContains := contains(allowedExt, filepath.Ext(icon.Filename))

	if !isContains {
		errorResponse(context, http.StatusUnprocessableEntity, "Icon must be jpg, jpeg or png")
		return
	}

	_ := h.storageProvider.ConnectAws()

	category, err := h.services.Categories.Create(context.Request.Context(), categoryDTO)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, category)
}
