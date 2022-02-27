package service

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"github.com/sigit14ap/go-commerce/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoresService struct {
	repo repository.Stores
}

func NewStoresService(repo repository.Stores) *StoresService {
	return &StoresService{
		repo: repo,
	}
}

func (service *StoresService) FindByUserID(ctx context.Context, userID primitive.ObjectID) (domain.Store, error) {
	return service.repo.FindByUserID(ctx, userID)
}

func (service *StoresService) Create(ctx context.Context, store dto.StoreRegisterDTO) (domain.Store, error) {
	return service.repo.Create(ctx, store)
}
