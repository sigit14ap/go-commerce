package service

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressesService struct {
	repo repository.Addresses
}

func (service *AddressesService) FindAll(ctx context.Context, userID primitive.ObjectID) ([]domain.Address, error) {
	return service.repo.FindAll(ctx, userID)
}

func NewAddressesService(repo repository.Addresses) *AddressesService {
	return &AddressesService{
		repo: repo,
	}
}
