package repository

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type AddressesRepo struct {
	db *mongo.Collection
}

func (repo AddressesRepo) FindAll(ctx context.Context, userID primitive.ObjectID) ([]domain.Address, error) {
	cursor, err := repo.db.Find(ctx, bson.M{"userID": userID})
	if err != nil {
		return nil, err
	}

	data := []domain.Address{}
	err = cursor.All(ctx, &data)
	return data, err
}

func NewAddressesRepo(db *mongo.Database) *AddressesRepo {
	collection := db.Collection(addressesCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"UserID": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("unable to create cart collection index, %v", err)
	}

	return &AddressesRepo{
		db: collection,
	}
}
