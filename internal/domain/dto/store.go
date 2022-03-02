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

type StoreShipmentInput struct {
	ProvinceID string `json:"province_id" validate:"required"`
	CityID     string `json:"city_id" validate:"required"`
}

type StoreShipmentDTO struct {
	ProvinceID primitive.ObjectID `json:"province_id" bson:"province_id"`
	CityID     primitive.ObjectID `json:"city_id" bson:"city_id"`
}
