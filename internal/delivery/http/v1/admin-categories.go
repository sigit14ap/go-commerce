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

// GetCategories godoc
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

// GetCategoryByIdAdmin godoc
// @Summary   Get category by id
// @Tags      admin-categories
// @Accept    json
// @Produce   json
// @Param     id   path      string  true  "category id"
// @Success   200  {object}  success
// @Failure   400  {object}  failure
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  AdminAuth
// @Router    /admins/categories/{id} [get]
func (h *Handler) getCategoryByIdAdmin(context *gin.Context) {

	categoryID, err := getIdFromPath(context, "id")
	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.services.Categories.FindByID(context.Request.Context(), categoryID)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, category)
}

// CreateCategory godoc
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
	var categoryInput dto.ValidationCategoryDTO
	err := context.ShouldBindWith(&categoryInput, binding.FormMultipart)

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

	file, _ := context.FormFile("icon")
	uploadedFile := h.storageProvider.Upload("Category", file)

	var categoryDTO dto.CreateCategoryDTO
	categoryDTO.Name = categoryInput.Name
	categoryDTO.Description = categoryInput.Description
	categoryDTO.Icon = uploadedFile

	category, err := h.services.Categories.Create(context.Request.Context(), categoryDTO)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, category)
}

// UpdateCategory godoc
// @Summary   Update category
// @Tags      admin-categories
// @Accept    json
// @Produce   json
// @Param     id       path      string                true  "category id"
// @Param     category  body      dto.ValidationUpdateCategoryDTO  true  "category update fields"
// @Success   200  {object}  success
// @Failure   400  {object}  failure
// @Failure   401  {object}  failure
// @Failure   422  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  AdminAuth
// @Router    /admins/categories/{id} [put]
func (h *Handler) updateCategoryAdmin(context *gin.Context) {
	var categoryInput dto.ValidationUpdateCategoryDTO
	err := context.ShouldBind(&categoryInput)

	if err != nil {
		errorResponse(context, http.StatusUnprocessableEntity, "invalid input body")
		return
	}

	categoryID, err := getIdFromPath(context, "id")
	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	var categoryDTO dto.UpdateCategoryDTO
	categoryDTO.Name = categoryInput.Name
	categoryDTO.Description = categoryInput.Description

	icon, err := context.FormFile("icon")

	if err == nil {
		allowedExt := []string{".jpg", ".png", "jpeg"}
		isContains := contains(allowedExt, filepath.Ext(icon.Filename))

		if !isContains {
			errorResponse(context, http.StatusUnprocessableEntity, "Icon must be jpg, jpeg or png")
			return
		}

		file, _ := context.FormFile("icon")
		uploadedFile := h.storageProvider.Upload("Category", file)

		categoryDTO.Icon = uploadedFile
	}

	category, err := h.services.Categories.Update(context.Request.Context(), categoryDTO, categoryID)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, category)
}

// DeleteCategory godoc
// @Summary   Delete category
// @Tags      admin-categories
// @Accept    json
// @Produce   json
// @Param     id       path      string                true  "category id"
// @Success   200  {object}  success
// @Failure   400  {object}  failure
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  AdminAuth
// @Router    /admins/categories/{id} [get]
func (h *Handler) deleteCategoryAdmin(context *gin.Context) {

	categoryID, err := getIdFromPath(context, "id")
	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Categories.Delete(context.Request.Context(), categoryID)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	var category interface{}
	successResponse(context, category)
}
