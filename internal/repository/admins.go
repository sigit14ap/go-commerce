package repository

import (
	"context"

	"github.com/sigit14ap/go-commerce/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminsRepo struct {
	db *mongo.Collection
}

func (a AdminsRepo) FindByCredentials(ctx context.Context, email string) (domain.Admin, error) {
	result := a.db.FindOne(ctx, bson.M{"email": email})

	var admin domain.Admin
	err := result.Decode(&admin)

	return admin, err
}

func NewAdminsRepo(db *mongo.Database) *AdminsRepo {
	return &AdminsRepo{
		db: db.Collection(adminsCollection),
	}
}
