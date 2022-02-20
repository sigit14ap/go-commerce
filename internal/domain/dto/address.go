package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type AddressDTO struct {
	UserID      primitive.ObjectID `json:"user_id"`
	Fullname    string             `json:"fullname"`
	PhoneNumber string             `json:"phone_number"`
	ProvinceID  primitive.ObjectID `json:"province_id"`
	CityID      primitive.ObjectID `json:"city_id"`
	Address     string             `json:"address"`
	Longitude   string             `json:"longitude"`
	Latitude    string             `json:"latitude"`
	Type        string             `json:"type"`
	IsPrimary   bool               `json:"is_primary"`
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
