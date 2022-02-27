package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"net/http"
)

func (h *Handler) initStoreRoutes(api *gin.RouterGroup) {
	store := api.Group("/store", h.verifyUser)
	{
		store.POST("/register", h.storeRegister)

	}
}

// StoreRegister godoc
// @Summary  Register store
// @Tags     store
// @Accept   json
// @Produce  json
// @Param    store  body      dto.StoreRegisterInput  true  "register"
// @Success  200       {object}  success
// @Failure  400       {object}  failure
// @Failure  401       {object}  failure
// @Failure  404       {object}  failure
// @Failure  500       {object}  failure
// @Router   /store/register [post]
func (h *Handler) storeRegister(context *gin.Context) {
	userID, err := getIdFromRequestContext(context, "userID")
	if err != nil {
		errorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	var input dto.StoreRegisterInput
	_ = context.ShouldBindJSON(&input)

	err = validate.Struct(input)
	if err != nil {
		errorValidationResponse(context, err)
		return
	}

	registerDTO := dto.StoreRegisterDTO{}
	copier.Copy(&registerDTO, &input)
	registerDTO.UserID = userID

	_, err = h.services.Stores.FindByUserID(context, userID)

	if err == nil {
		errorResponse(context, http.StatusForbidden, "Already registered as store")
		return
	}

	data, err := h.services.Stores.Create(context, registerDTO)

	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, data)
}
