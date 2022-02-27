package repository

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
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
	cursor, err := repo.db.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}

	data := []domain.Address{}
	err = cursor.All(ctx, &data)
	return data, err
}

func (repo *AddressesRepo) Find(ctx context.Context, userID primitive.ObjectID, addressID primitive.ObjectID) (domain.Address, error) {
	result := repo.db.FindOne(ctx, bson.M{"_id": addressID, "userId": userID})

	var address domain.Address
	err := result.Decode(&address)

	return address, err
}

func (repo *AddressesRepo) Create(ctx context.Context, address dto.AddressDTO) (domain.Address, error) {
	result, err := repo.db.InsertOne(ctx, address)

	data := domain.Address{}
	copier.Copy(&data, &address)
	data.ID = result.InsertedID.(primitive.ObjectID)

	return data, err
}

func (repo *AddressesRepo) Update(ctx context.Context, userID primitive.ObjectID, addressID primitive.ObjectID, address dto.AddressDTO) (domain.Address, error) {
	_, err := repo.db.UpdateOne(ctx, bson.M{"_id": addressID, "userId": userID}, bson.M{"$set": address})

	data := domain.Address{}
	copier.Copy(&data, &address)
	data.ID = addressID

	return data, err
}

func (repo *AddressesRepo) Delete(ctx context.Context, userID primitive.ObjectID, addressID primitive.ObjectID) error {
	_, err := repo.db.DeleteOne(ctx, bson.M{"_id": addressID, "userId": userID})
	return err
}

func NewAddressesRepo(db *mongo.Database) *AddressesRepo {
	collection := db.Collection(addressesCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"UserID": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("unable to address cart collection index, %v", err)
	}

	return &AddressesRepo{
		db: collection,
	}
}
