package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type StoreRegisterInput struct {
	Name   string `json:"name" validate:"required,min=5,max=255"`
	Domain string `json:"domain" validate:"required,alphanum,min=5,max=20"`
}

type StoreRegisterDTO struct {
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`
	Name   string             `json:"name" bson:"name"`
	Domain string             `json:"domain" bson:"domain"`
}
