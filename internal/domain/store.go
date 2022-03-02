package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Store struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name               string             `json:"name" bson:"name"`
	Domain             string             `json:"domain" bson:"domain"`
	ShipmentCityID     primitive.ObjectID `json:"shipment_city_id" bson:"shipment_city_id"`
	ShipmentProvinceID primitive.ObjectID `json:"shipment_province_id" bson:"shipment_province_id"`
}
