package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Province struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ThirdPartyID string             `json:"third_party_id" bson:"third_party_id"`
	Name         string             `json:"name"`
}

type City struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProvinceID   primitive.ObjectID `json:"province_id" bson:"province_id"`
	ThirdPartyID string             `json:"third_party_id" bson:"third_party_id"`
	Name         string             `json:"name" bson:"name"`
}
