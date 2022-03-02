package repository

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductsRepo struct {
	db *mongo.Collection
}

func (p ProductsRepo) GetBySellerID(ctx context.Context, storeID primitive.ObjectID) ([]domain.Product, error) {
	cursor, err := p.db.Find(ctx, bson.M{"store_id": storeID})
	if err != nil {
		return nil, err
	}

	var productArray []domain.Product
	err = cursor.All(ctx, &productArray)
	return productArray, err
}

func (p ProductsRepo) FindAll(ctx context.Context) ([]domain.Product, error) {
	cursor, err := p.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var productArray []domain.Product
	err = cursor.All(ctx, &productArray)
	return productArray, err
}

func (p ProductsRepo) FindByID(ctx context.Context, productID primitive.ObjectID) (domain.Product, error) {
	result := p.db.FindOne(ctx, bson.M{"_id": productID})

	var product domain.Product
	err := result.Decode(&product)

	return product, err
}

func (p ProductsRepo) Create(ctx context.Context, product domain.Product) (domain.Product, error) {

	images := []domain.ProductImage{}

	for _, image := range product.Images {
		images = append(images, domain.ProductImage{
			ID:    primitive.NewObjectID(),
			Image: image.Image,
		})
	}

	product.ID = primitive.NewObjectID()
	product.Images = images
	_, err := p.db.InsertOne(ctx, product)
	return product, err
}

func (p ProductsRepo) Update(ctx context.Context, product domain.Product, productID primitive.ObjectID) (domain.Product, error) {
	updateQuery := bson.M{}

	if product.Name != "" {
		updateQuery["name"] = product.Name
	}

	if product.Description != "" {
		updateQuery["description"] = product.Description
	}

	_, err := p.db.UpdateOne(ctx, bson.M{"_id": productID}, bson.M{"$set": updateQuery})
	findResult := p.db.FindOne(ctx, bson.M{"_id": productID})

	var result domain.Product
	err = findResult.Decode(&result)

	return result, err
}

func (p ProductsRepo) Delete(ctx context.Context, productID primitive.ObjectID, storeID primitive.ObjectID) error {
	_, err := p.db.DeleteOne(ctx, bson.M{"_id": productID, "storeID": storeID})
	return err
}

func NewProductsRepo(db *mongo.Database) *ProductsRepo {
	return &ProductsRepo{
		db: db.Collection(productsCollection),
	}
}
