package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"net/http"
)

func (h *Handler) initUserAddressRoutes(api *gin.RouterGroup) {
	address := api.Group("/users/address", h.verifyUser)
	{
		address.GET("/", h.getAddress)
		address.POST("/", h.createAddress)
		address.PUT("/:addressID", h.updateAddress)
		address.DELETE("/:addressID", h.deleteAddress)
	}
}

// GetAddress godoc
// @Summary  Get user address
// @Tags     address
// @Accept   json
// @Produce  json
// @Success  200     {array}   success
// @Failure  401     {object}  failure
// @Failure  404     {object}  failure
// @Failure  500     {object}  failure
// @Router   /user/address [get]
func (h *Handler) getAddress(context *gin.Context) {
	userID, err := getIdFromRequestContext(context, "userID")
	if err != nil {
		errorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	data, err := h.services.Addresses.FindAll(context, userID)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, data)
}

// CreateAddress godoc
// @Summary  Add address
// @Tags     address
// @Accept   json
// @Produce  json
// @Param    address  body      dto.AddressInput  true  "address"
// @Success  201       {object}  success
// @Failure  400       {object}  failure
// @Failure  401       {object}  failure
// @Failure  404       {object}  failure
// @Failure  500       {object}  failure
// @Router   /users/address [post]
func (h *Handler) createAddress(context *gin.Context) {
	userID, err := getIdFromRequestContext(context, "userID")
	if err != nil {
		errorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	var input dto.AddressInput
	_ = context.ShouldBindJSON(&input)

	err = validate.Struct(input)
	if err != nil {
		errorValidationResponse(context, err)
		return
	}

	provinceID, err := getIdFromRequest(input.ProvinceID)

	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	cityID, err := getIdFromRequest(input.CityID)

	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	addressDTO := dto.AddressDTO{}
	copier.Copy(&addressDTO, &input)
	addressDTO.UserID = userID
	addressDTO.ProvinceID = provinceID
	addressDTO.CityID = cityID

	_, err = h.services.Areas.FindCityAndProvince(context, cityID, provinceID)

	if err != nil {
		errorResponse(context, http.StatusBadRequest, "City not found")
		return
	}

	data, err := h.services.Addresses.Create(context, addressDTO)

	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	createdResponse(context, data)
}

// UpdateAddress godoc
// @Summary  Update address
// @Tags     address
// @Accept   json
// @Produce  json
// @Param    addressID  path      string  true  "address id"
// @Param    address  	body      dto.AddressInput  true  "address"
// @Success  201       {object}  success
// @Failure  400       {object}  failure
// @Failure  401       {object}  failure
// @Failure  404       {object}  failure
// @Failure  500       {object}  failure
// @Router   /users/address/{addressID} [put]
func (h *Handler) updateAddress(context *gin.Context) {
	userID, err := getIdFromRequestContext(context, "userID")
	if err != nil {
		errorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	addressID, err := getIdFromPath(context, "addressID")
	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.services.Addresses.Find(context, userID, addressID)

	if err != nil {
		errorResponse(context, http.StatusNotFound, "Address not found")
		return
	}

	var input dto.AddressInput
	_ = context.ShouldBindJSON(&input)

	err = validate.Struct(input)
	if err != nil {
		errorValidationResponse(context, err)
		return
	}

	provinceID, err := getIdFromRequest(input.ProvinceID)

	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	cityID, err := getIdFromRequest(input.CityID)

	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	addressDTO := dto.AddressDTO{}
	copier.Copy(&addressDTO, &input)
	addressDTO.UserID = userID
	addressDTO.ProvinceID = provinceID
	addressDTO.CityID = cityID

	_, err = h.services.Areas.FindCityAndProvince(context, cityID, provinceID)

	if err != nil {
		errorResponse(context, http.StatusBadRequest, "City not found")
		return
	}

	data, err := h.services.Addresses.Update(context, userID, addressID, addressDTO)

	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, data)
}

// DeleteAddress godoc
// @Summary  Delete address
// @Tags     address
// @Accept   json
// @Produce  json
// @Param    addressID  path      string  true  "address id"
// @Success  201       {object}  success
// @Failure  400       {object}  failure
// @Failure  401       {object}  failure
// @Failure  404       {object}  failure
// @Failure  500       {object}  failure
// @Router   /users/address/{addressID} [delete]
func (h *Handler) deleteAddress(context *gin.Context) {
	userID, err := getIdFromRequestContext(context, "userID")
	if err != nil {
		errorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	addressID, err := getIdFromPath(context, "addressID")
	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.services.Addresses.Find(context, userID, addressID)

	if err != nil {
		errorResponse(context, http.StatusNotFound, "Address not found")
		return
	}

	err = h.services.Addresses.Delete(context, userID, addressID)

	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	var data interface{}
	successResponse(context, data)
}
