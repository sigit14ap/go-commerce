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

func (service *StoresService) FindByDomain(ctx context.Context, domainStore string) (domain.Store, error) {
	return service.repo.FindByDomain(ctx, domainStore)
}

func (service *StoresService) Create(ctx context.Context, store dto.StoreRegisterDTO) (domain.Store, error) {
	return service.repo.Create(ctx, store)
}

func (service *StoresService) UpdateShipment(ctx context.Context, storeID primitive.ObjectID, shipment dto.StoreShipmentDTO) (domain.Store, error) {
	return service.repo.UpdateShipment(ctx, storeID, shipment)
}
