package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Store struct {
	ID     primitive.ObjectID `json:"id"`
	Name   string             `json:"name"`
	Domain string             `json:"domain"`
}
