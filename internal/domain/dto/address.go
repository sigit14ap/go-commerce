package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type AddressDTO struct {
	UserID      primitive.ObjectID `json:"user_id" bson:"userId"`
	Fullname    string             `json:"fullname" bson:"fullname"`
	PhoneNumber string             `json:"phone_number" bson:"phoneNumber"`
	ProvinceID  primitive.ObjectID `json:"province_id" bson:"provinceId"`
	CityID      primitive.ObjectID `json:"city_id" bson:"cityId"`
	Address     string             `json:"address" bson:"address"`
	Longitude   string             `json:"longitude" bson:"longitude"`
	Latitude    string             `json:"latitude" bson:"latitude"`
	Type        string             `json:"type" bson:"type"`
	IsPrimary   bool               `json:"is_primary" bson:"isPrimary"`
}

type AddressInput struct {
	Fullname    string `json:"fullname" validate:"required,min=5,max=255"`
	PhoneNumber string `json:"phone_number" validate:"required,numeric,min=10,max=15"`
	ProvinceID  string `json:"province_id" validate:"required"`
	CityID      string `json:"city_id" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Longitude   string `json:"longitude" validate:"required"`
	Latitude    string `json:"latitude" validate:"required"`
	Type        string `json:"type" validate:"required,oneof='Office' 'Home'"`
	IsPrimary   *bool  `json:"is_primary" validate:"required"`
}
