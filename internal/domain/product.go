package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	TotalRating float64            `json:"total_rating" bson:"-"`
	CategoryID  string             `json:"-" bson:"category_id"`
	Category    Category           `json:"category" bson:"-"`
	Images      []ProductImage     `json:"images" bson:"images"`
}

type ProductImage struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Image string             `json:"image" bson:"image"`
}

type Category struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Icon        string             `json:"icon" bson:"icon"`
}
