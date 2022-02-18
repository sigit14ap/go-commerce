package dto

import (
	"github.com/sigit14ap/go-commerce/internal/domain"
)

type CreateCartDTO struct {
	CartItems []domain.CartItem `json:"products"`
}

type UpdateCartDTO struct {
	CartItems []domain.CartItem `json:"products"`
}

type UpdateCartInput struct {
	CartItems []domain.CartItem `json:"products"`
}

type UpdateCartItemDTO struct {
	Quantity int64 `json:"quantity"`
}

type AddToCartDTO struct {
	ProductID string `json:"productID" bson:"productID"`
	Quantity  int64  `json:"quantity" bson:"quantity"`
}
