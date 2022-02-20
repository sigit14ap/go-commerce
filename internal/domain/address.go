package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Fullname    string             `json:"fullname" bson:"fullname"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	ProvinceID  string             `json:"province_id" bson:"province_id"`
	CityID      string             `json:"city_id" bson:"city_id"`
	Address     string             `json:"address" bson:"address"`
	Longitude   string             `json:"longitude" bson:"longitude"`
	Latitude    string             `json:"latitude" bson:"latitude"`
	Type        string             `json:"type" bson:"type"`
	IsPrimary   bool               `json:"is_primary" bson:"is_primary"`
}
