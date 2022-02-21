package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"net/http"
)

func (h *Handler) initAreasRoutes(api *gin.RouterGroup) {
	area := api.Group("/area")
	{
		area.GET("/province", h.getProvince)
		area.GET("/city", h.getCity)
	}
}

// GetProvince godoc
// @Summary   Get all province
// @Tags      area-province
// @Accept    json
// @Produce   json
// @Success   200  {array}   success
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Router    /area/province [get]
func (h *Handler) getProvince(context *gin.Context) {
	provinces, err := h.services.Areas.GetProvinces(context)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, provinces)
}

// GetCity godoc
// @Summary   Get all city
// @Tags      area-city
// @Accept    json
// @Produce   json
// @Param    city  body      dto.CityListDTO
// @Success   200  {array}   success
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Router    /area/city [get]
func (h *Handler) getCity(context *gin.Context) {
	var cityListInput dto.CityListInput

	err := context.BindJSON(&cityListInput)
	if err != nil {
		errorResponse(context, http.StatusBadRequest, "invalid input body")
		return
	}

	provinceID, err := getIdFromRequest(cityListInput.ProvinceID)
	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	cityListDTO := dto.CityListDTO{
		ProvinceID: provinceID,
	}

	cities, err := h.services.Areas.GetCities(context, cityListDTO)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, cities)
}
