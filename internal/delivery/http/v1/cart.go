package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func (h *Handler) initCartRoutes(api *gin.RouterGroup) {
	cart := api.Group("/cart", h.verifyUser)
	{
		cart.GET("/", h.getCartItems)
		cart.POST("/", h.createCartItem)
		cart.DELETE("/", h.clearCart)
		cart.PUT("/:productID", h.updateCartItem)
		cart.DELETE("/:productID", h.deleteCartItem)
	}
}

// GetCartItems godoc
// @Summary  Get cart items
// @Tags     cart
// @Accept   json
// @Produce  json
// @Param    Cookie  header    string  true  "cart id"
// @Success  200     {array}   success
// @Failure  401     {object}  failure
// @Failure  404     {object}  failure
// @Failure  500     {object}  failure
// @Router   /cart [get]
func (h *Handler) getCartItems(context *gin.Context) {
	userID, err := getIdFromRequestContext(context, "userID")
	if err != nil {
		errorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	cartItems, err := h.services.Carts.FindCartItems(context, userID)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, cartItems)
}

// CreateCartItem godoc
// @Summary  Add cart item
// @Tags     cart
// @Accept   json
// @Produce  json
// @Param    cartItem  body      domain.CartItem  true  "cart item"
// @Param    Cookie    header    string           true  "cart id"
// @Success  201       {object}  success
// @Failure  400       {object}  failure
// @Failure  401       {object}  failure
// @Failure  404       {object}  failure
// @Failure  500       {object}  failure
// @Router   /cart [post]
func (h *Handler) createCartItem(context *gin.Context) {
	userID, err := getIdFromRequestContext(context, "userID")
	if err != nil {
		errorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	var cartItemInput dto.AddToCartDTO
	err = context.BindJSON(&cartItemInput)
	if err != nil {
		errorResponse(context, http.StatusBadRequest, "invalid input body")
		return
	}

	productID, err := getIdFromRequest(cartItemInput.ProductID)
	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	cartData := domain.CartItem{
		ProductID: productID,
		Quantity:  cartItemInput.Quantity,
	}

	cartItem, err := h.services.Carts.AddCartItem(context, cartData, userID)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	createdResponse(context, cartItem)
}

// UpdateCartItem godoc
// @Summary  Update cart item
// @Tags     cart
// @Accept   json
// @Produce  json
// @Param    productID  path      string                 true  "product id"
// @Param    cartItem   body      dto.UpdateCartItemDTO  true  "cart item"
// @Param    Cookie     header    string                 true  "cart id"
// @Success  200     {object}  success
// @Failure  400     {object}  failure
// @Failure  401     {object}  failure
// @Failure  404     {object}  failure
// @Failure  500     {object}  failure
// @Router   /cart/{productID} [put]
func (h *Handler) updateCartItem(context *gin.Context) {
	userID, err := getIdFromRequestContext(context, "userID")
	if err != nil {
		errorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	productID, err := getIdFromPath(context, "productID")
	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	var cartItemInput dto.UpdateCartItemDTO
	err = context.BindJSON(&cartItemInput)
	if err != nil {
		errorResponse(context, http.StatusBadRequest, "invalid input body")
		return
	}

	cartItem, err := h.services.Carts.UpdateCartItem(context, domain.CartItem{
		ProductID: productID,
		Quantity:  cartItemInput.Quantity,
	}, userID)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, cartItem)
}

// DeleteCartItem godoc
// @Summary  Delete cart item
// @Tags     cart
// @Accept   json
// @Produce  json
// @Param    productID  path      string  true  "product id"
// @Param    Cookie  header    string  true  "cart id"
// @Success  200        {object}  success
// @Failure  400        {object}  failure
// @Failure  401        {object}  failure
// @Failure  404        {object}  failure
// @Failure  500        {object}  failure
// @Router   /cart/{productID} [delete]
func (h *Handler) deleteCartItem(context *gin.Context) {
	cartIDHex, ok := context.Get("cartID")
	if !ok {
		errorResponse(context, http.StatusInternalServerError, "failed to get cart id")
		return
	}
	cartID := cartIDHex.(primitive.ObjectID)

	productID, err := getIdFromPath(context, "productID")
	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Carts.DeleteCartItem(context, productID, cartID)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	context.Status(http.StatusOK)
}

// ClearCart godoc
// @Summary  Delete all cart items
// @Tags     cart
// @Accept   json
// @Produce  json
// @Param    Cookie     header    string  true  "cart id"
// @Success  200        {object}  success
// @Failure  400        {object}  failure
// @Failure  401        {object}  failure
// @Failure  404        {object}  failure
// @Failure  500        {object}  failure
// @Router   /cart [delete]
func (h *Handler) clearCart(context *gin.Context) {
	cartIDHex, ok := context.Get("cartID")
	if !ok {
		errorResponse(context, http.StatusInternalServerError, "failed to get cart id")
		return
	}
	cartID := cartIDHex.(primitive.ObjectID)

	err := h.services.Carts.ClearCart(context, cartID)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	context.Status(http.StatusOK)
}

// GetCarts godoc
// @Summary   Get all carts
// @Tags      admin-carts
// @Accept    json
// @Produce   json
// @Success   200  {array}   success
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  AdminAuth
// @Router    /admins/carts [get]
func (h *Handler) getAllCartsAdmin(context *gin.Context) {
	carts, err := h.services.Carts.FindAll(context.Request.Context())
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	cartsArray := make([]domain.Cart, len(carts))
	if carts != nil {
		cartsArray = carts
	}

	successResponse(context, cartsArray)
}

// GetCartByIdAdmin godoc
// @Summary   Get cart by id
// @Tags      admin-carts
// @Accept    json
// @Produce   json
// @Param     id   path      string  true  "cart id"
// @Success   200  {object}  success
// @Failure   400  {object}  failure
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  AdminAuth
// @Router    /admins/carts/{id} [get]
func (h *Handler) getCartByIdAdmin(context *gin.Context) {
	id, err := getIdFromPath(context, "id")
	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}
	cart, err := h.services.Carts.FindByID(context.Request.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			errorResponse(context, http.StatusInternalServerError,
				fmt.Sprintf("no carts with id: %s", id.Hex()))
		} else {
			errorResponse(context, http.StatusInternalServerError, err.Error())
		}
		return
	}

	successResponse(context, cart)
}

// DeleteCart godoc
// @Summary   Delete cart
// @Tags      admin-carts
// @Accept    json
// @Produce   json
// @Param     id   path      string  true  "cart id"
// @Success   200  {object}  success
// @Failure   400  {object}  failure
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  AdminAuth
// @Router    /admins/carts/{id} [delete]
func (h *Handler) deleteCartAdmin(context *gin.Context) {
	cartID, err := getIdFromPath(context, "id")
	if err != nil {
		errorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Carts.Delete(context, cartID)
	if err != nil {
		errorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	context.Status(http.StatusOK)
}
