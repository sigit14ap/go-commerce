package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type CityListInput struct {
	ProvinceID string `json:"province_id"`
}

type CityListDTO struct {
	ProvinceID primitive.ObjectID `json:"province_id"`
}

type CreateCityDTO struct {
	ProvinceID primitive.ObjectID `json:"province_id"`
	Name       string             `json:"name"`
}

type ThirdPartyCityDTO struct {
	ProvinceID string `json:"province_id"`
	CityID     string `json:"city_id"`
	Name       string `json:"name"`
}
