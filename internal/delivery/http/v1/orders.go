package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
)

func (h *Handler) initOrdersRoutes(api *gin.RouterGroup) {
	orders := api.Group("/users/orders", h.verifyUser)
	{
		orders.GET("/delivery-cost", h.getDeliveryCost)
		orders.GET("/", h.getUserOrders)
		orders.POST("/", h.createOrder)
		orders.GET("/payment/:id", h.getOrderPaymentLink)
	}
}

// GetDeliveryCost godoc
// @Summary   Delivery cost List
// @Tags      user
// @Accept    json
// @Produce   json
// @Success   200  {array}   success
// @Failure   401    {object}  failure
// @Failure   404    {object}  failure
// @Failure   500    {object}  failure
// @Security  UserAuth
// @Router    /users/orders/delivery-cost [get]
func (h *Handler) getDeliveryCost(context *gin.Context) {
	userID, err := getIdFromRequestContext(context, "userID")
	if err != nil {
		ErrorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	var input dto.DeliveryCostInput
	_ = context.ShouldBindJSON(&input)

	err = validate.Struct(input)
	if err != nil {
		errorValidationResponse(context, err)
		return
	}

	for _, product := range input.Product {
		productID, err := getIdFromRequest(product.ProductID)

		if err != nil {
			ErrorResponse(context, http.StatusBadRequest, err.Error())
			return
		}

		_, err = h.services.Carts.FindItem(context, userID, productID)

		if err != nil {
			ErrorResponse(context, http.StatusBadRequest, "Product id "+product.ProductID+" not found in cart")
			return
		}
	}

	successResponse(context, input)
}

// GerUserOrders godoc
// @Summary   User order List
// @Tags      user
// @Accept    json
// @Produce   json
// @Success   200  {array}   success
// @Failure   401    {object}  failure
// @Failure   404    {object}  failure
// @Failure   500    {object}  failure
// @Security  UserAuth
// @Router    /users/orders [get]
func (h *Handler) getUserOrders(context *gin.Context) {
	userID, err := getIdFromRequestContext(context, "userID")
	if err != nil {
		ErrorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	orders, err := h.services.Orders.FindByUserID(context.Request.Context(), userID)
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, orders)
}

// CreateOrder godoc
// @Summary   Create order
// @Tags      user
// @Accept    json
// @Produce   json
// @Param     order  body      dto.CreateOrderDTO  true  "contact info"
// @Success   201    {object}  success
// @Failure   400  {object}  failure
// @Failure   401    {object}  failure
// @Failure   404    {object}  failure
// @Failure   500    {object}  failure
// @Security  UserAuth
// @Router    /users/orders [post]
func (h *Handler) createOrder(context *gin.Context) {
	userID, err := getIdFromRequestContext(context, "userID")
	if err != nil {
		ErrorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	cart, err := h.services.Carts.FindByID(context.Request.Context(), userID)
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	if len(cart.CartItems) == 0 {
		ErrorResponse(context, http.StatusBadRequest, "user cart is empty")
		return
	}

	orderItems := make([]domain.OrderItem, len(cart.CartItems))
	for i, cartItem := range cart.CartItems {
		orderItems[i] = domain.OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
		}
	}

	var createOrderDTO dto.CreateOrderDTO
	err = context.BindJSON(&createOrderDTO)
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, "invalid input body")
		return
	}

	order, err := h.services.Orders.Create(context.Request.Context(), dto.CreateOrderDTO{
		OrderItems:  orderItems,
		ContactInfo: createOrderDTO.ContactInfo,
		UserID:      userID,
	})

	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Carts.ClearCart(context.Request.Context(), cart.ID)
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, "cart can't be cleared")
		return
	}

	successResponse(context, order)
}

// PaymentLink godoc
// @Summary   Get order payment link
// @Tags      user
// @Accept    json
// @Produce   json
// @Param     id   path      string  true  "order id"
// @Success   200  {object}  success
// @Failure   400    {object}  failure
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  UserAuth
// @Router    /users/orders/{id}/payment [get]
func (h *Handler) getOrderPaymentLink(context *gin.Context) {
	orderID, err := getIdFromPath(context, "id")
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	link, err := h.services.Payment.GetPaymentLink(context.Request.Context(), orderID)
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, link)
}

// GetOrdersAdmin godoc
// @Summary   Get all orders
// @Tags      admin-orders
// @Accept    json
// @Produce   json
// @Success   200  {array}   success
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  AdminAuth
// @Router    /admins/orders [get]
func (h *Handler) getAllOrdersAdmin(context *gin.Context) {
	orders, err := h.services.Orders.FindAll(context.Request.Context())
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	orderArray := make([]domain.Order, len(orders))
	if orders != nil {
		orderArray = orders
	}

	successResponse(context, orderArray)
}

// UpdateOrder godoc
// @Summary   Update order
// @Tags      admin-orders
// @Accept    json
// @Produce   json
// @Param     id     path      string              true  "order id"
// @Param     order  body      dto.UpdateOrderDTO  true  "order update fields"
// @Success   200    {object}  success
// @Failure   400    {object}  failure
// @Failure   401  {object}  failure
// @Failure   404  {object}  failure
// @Failure   500  {object}  failure
// @Security  AdminAuth
// @Router    /admins/orders/{id} [put]
func (h *Handler) updateOrderAdmin(context *gin.Context) {
	var orderDTO dto.UpdateOrderDTO

	err := context.BindJSON(&orderDTO)
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, "invalid input body")
		return
	}

	orderID, err := getIdFromPath(context, "id")
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.services.Orders.Update(context.Request.Context(), orderDTO, orderID)
	if err != nil {
		ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(context, order)
}
