package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
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
	Province    Province           `json:"province"`
	City        City               `json:"city"`
}
