package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"net/http"
)

func (h *Handler) initAreasRoutes(api *gin.RouterGroup) {
	users := api.Group("/area")
	{
		users.GET("/province", h.getProvince)
		users.GET("/city", h.getCity)
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
	var cityListDTO dto.CityListDTO

	err := context.BindJSON(&cityListDTO)
	if err != nil {
		errorResponse(context, http.StatusBadRequest, "invalid input body")
		return
	}

	cities, err := h.services.Areas.GetCities(context, cityListDTO)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, cities)
}
