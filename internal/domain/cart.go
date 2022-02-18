package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	UserID     primitive.ObjectID `json:"userID" bson:"userID"`
	TotalPrice float64            `json:"totalPrice" bson:"-"`
	CartItems  []CartItem         `json:"cartItems" bson:"cartItems"`
}

type CartItem struct {
	ProductID primitive.ObjectID `json:"productID" bson:"productID"`
	Quantity  int64              `json:"quantity" bson:"quantity"`
	Product   Product            `json:"product" bson:"-"`
}
