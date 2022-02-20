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
// @Param    address  body      dto.AddressInput  true  "cart item"
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

	addressDTO := dto.AddressDTO{}
	copier.Copy(&addressDTO, &input)
	addressDTO.UserID = userID

	//cartItem, err := h.services.Carts.AddCartItem(context, cartData, userID)
	//if err != nil {
	//	errorResponse(context, http.StatusInternalServerError, err.Error())
	//	return
	//}

	createdResponse(context, addressDTO)
}
