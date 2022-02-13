package repository

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoriesRepo struct {
	db *mongo.Collection
}

func (repo CategoriesRepo) FindAll(ctx context.Context) ([]domain.Category, error) {
	cursor, err := repo.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var categoryArray []domain.Category
	err = cursor.All(ctx, &categoryArray)
	return categoryArray, err
}

func (repo CategoriesRepo) FindByID(ctx context.Context, categoryID primitive.ObjectID) (domain.Category, error) {
	result := repo.db.FindOne(ctx, bson.M{"_id": categoryID})

	var category domain.Category
	err := result.Decode(&category)

	return category, err
}

func (repo CategoriesRepo) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	category.ID = primitive.NewObjectID()
	_, err := repo.db.InsertOne(ctx, category)
	return category, err
}

func (repo CategoriesRepo) Update(ctx context.Context, categoryInput dto.UpdateCategoryInput, categoryID primitive.ObjectID) (domain.Category, error) {
	updateQuery := bson.M{}

	if categoryInput.Name != "" {
		updateQuery["name"] = categoryInput.Name
	}

	if categoryInput.Description != nil {
		updateQuery["description"] = categoryInput.Description
	}

	_, err := repo.db.UpdateOne(ctx, bson.M{"_id": categoryID}, bson.M{"$set": updateQuery})
	findResult := repo.db.FindOne(ctx, bson.M{"_id": categoryID})

	var category domain.Category
	err = findResult.Decode(&category)

	return category, err
}

func (repo CategoriesRepo) Delete(ctx context.Context, categoryID primitive.ObjectID) error {
	_, err := repo.db.DeleteOne(ctx, bson.M{"_id": categoryID})
	return err
}

func NewCategoriesRepo(db *mongo.Database) *CategoriesRepo {
	return &CategoriesRepo{
		db: db.Collection(categoriesCollection),
	}
}
