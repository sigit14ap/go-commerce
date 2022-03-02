package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/internal/delivery/http/services"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"net/http"
)

func (h *Handler) initStoreSettingRoutes(api *gin.RouterGroup) {
	settings := api.Group("/settings")
	{
		settings.POST("/shipment", h.storeSettingShipment)
	}
}

// StoreSettingShipment godoc
// @Summary   Setting shipment store
// @Tags      store-setting
// @Accept    json
// @Produce   json
// @Param     shipment  body      dto.StoreShipmentDTO  true  "shipment"
// @Success   200  {object}  success
// @Failure   400  {object}  failure
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  StoreAuth
// @Router    /store/settings/shipment [post]
func (h *Handler) storeSettingShipment(context *gin.Context) {

	storeID, err := services.GetIdFromRequestContext(context, "storeID")

	if err != nil {
		services.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	var input dto.StoreShipmentInput
	_ = context.ShouldBindJSON(&input)

	err = validate.Struct(input)
	if err != nil {
		services.ErrorValidationResponse(context, err)
		return
	}

	provinceID, err := services.GetIdFromRequest(input.ProvinceID)

	if err != nil {
		services.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	cityID, err := services.GetIdFromRequest(input.CityID)

	if err != nil {
		services.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.services.Areas.FindCityAndProvince(context, cityID, provinceID)

	if err != nil {
		services.ErrorResponse(context, http.StatusBadRequest, "City not found")
		return
	}

	shipmentDTO := dto.StoreShipmentDTO{}
	shipmentDTO.CityID = cityID
	shipmentDTO.ProvinceID = provinceID

	store, err := h.services.Stores.UpdateShipment(context.Request.Context(), storeID, shipmentDTO)

	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	services.SuccessResponse(context, store)
}
