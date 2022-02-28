package repository

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StoresRepo struct {
	db *mongo.Collection
}

func (repo *StoresRepo) FindByUserID(ctx context.Context, userID primitive.ObjectID) (domain.Store, error) {
	result := repo.db.FindOne(ctx, bson.M{"user_id": userID})

	var store domain.Store
	err := result.Decode(&store)

	return store, err
}

func (repo *StoresRepo) FindByDomain(ctx context.Context, domainStore string) (domain.Store, error) {
	result := repo.db.FindOne(ctx, bson.M{"domain": domainStore})

	var store domain.Store
	err := result.Decode(&store)

	return store, err
}

func (repo *StoresRepo) Create(ctx context.Context, store dto.StoreRegisterDTO) (domain.Store, error) {
	result, err := repo.db.InsertOne(ctx, store)

	data := domain.Store{}
	copier.Copy(&data, &store)
	data.ID = result.InsertedID.(primitive.ObjectID)

	return data, err
}

func NewStoresRepo(db *mongo.Database) *StoresRepo {
	collection := db.Collection(storesCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"domain": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("unable to create store collection index, %v", err)
	}

	return &StoresRepo{
		db: collection,
	}
}
