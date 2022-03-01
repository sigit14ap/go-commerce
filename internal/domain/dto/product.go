package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateProductDTO struct {
	StoreID     primitive.ObjectID `form:"seller_id" bson:"seller_id"`
	Name        string             `form:"name" bson:"name"`
	Description string             `form:"description" bson:"description"`
	Price       float64            `form:"price" bson:"price"`
	CategoryID  primitive.ObjectID `form:"category_id" bson:"category_id"`
	Images      []string           `form:"images"`
	Weight      int64              `form:"weight" bson:"weight"`
}

type CreateProductInput struct {
	Name        string  `form:"name" binding:"required"`
	Description string  `form:"description" binding:"required"`
	Price       float64 `form:"price" binding:"required"`
	CategoryID  string  `form:"category_id" binding:"required"`
	Weight      int64   `form:"weight" binding:"required"`
}

type UpdateProductDTO struct {
	StoreID     primitive.ObjectID `form:"seller_id" bson:"seller_id"`
	Name        string             `form:"name" bson:"name"`
	Description string             `form:"description" bson:"description"`
	Price       float64            `form:"price" bson:"price"`
	CategoryID  primitive.ObjectID `form:"category_id" bson:"category_id"`
	Images      []string           `form:"images"`
	Weight      int64              `form:"weight" bson:"weight"`
}

type UpdateProductInput struct {
	Name        string   `form:"name" binding:"required"`
	Description string   `form:"description" binding:"required"`
	Price       float64  `form:"price" binding:"required"`
	CategoryID  string   `form:"category_id" binding:"required"`
	Images      []string `form:"images"`
	Weight      int64    `form:"weight" binding:"required"`
}
