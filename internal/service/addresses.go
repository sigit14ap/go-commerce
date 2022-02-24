package service

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"github.com/sigit14ap/go-commerce/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressesService struct {
	area Areas
	repo repository.Addresses
}

func (service *AddressesService) FindAll(ctx context.Context, userID primitive.ObjectID) ([]domain.Address, error) {
	data, err := service.repo.FindAll(ctx, userID)

	if err != nil {
		return []domain.Address{}, err
	}

	for i, address := range data {
		data[i].Province, err = service.area.FindProvince(ctx, address.ProvinceID)

		if err != nil {
			return []domain.Address{}, err
		}

		data[i].City, err = service.area.FindCity(ctx, address.CityID)

		if err != nil {
			return []domain.Address{}, err
		}
	}

	return data, nil
}

func (service *AddressesService) Find(ctx context.Context, userID primitive.ObjectID, addressID primitive.ObjectID) (domain.Address, error) {
	data, err := service.repo.Find(ctx, userID, addressID)

	if err != nil {
		return domain.Address{}, err
	}

	data.Province, err = service.area.FindProvince(ctx, data.ProvinceID)

	if err != nil {
		return domain.Address{}, err
	}

	data.City, err = service.area.FindCity(ctx, data.CityID)

	if err != nil {
		return domain.Address{}, err
	}

	return data, nil
}

func (service *AddressesService) Create(ctx context.Context, address dto.AddressDTO) (domain.Address, error) {
	data, err := service.repo.Create(ctx, address)

	if err != nil {
		return domain.Address{}, err
	}

	data.Province, err = service.area.FindProvince(ctx, data.ProvinceID)

	if err != nil {
		return domain.Address{}, err
	}

	data.City, err = service.area.FindCity(ctx, data.CityID)

	if err != nil {
		return domain.Address{}, err
	}

	return data, nil
}

func (service *AddressesService) Update(ctx context.Context, userID primitive.ObjectID, addressID primitive.ObjectID, address dto.AddressDTO) (domain.Address, error) {
	data, err := service.repo.Update(ctx, userID, addressID, address)

	if err != nil {
		return domain.Address{}, err
	}

	data.Province, err = service.area.FindProvince(ctx, data.ProvinceID)

	if err != nil {
		return domain.Address{}, err
	}

	data.City, err = service.area.FindCity(ctx, data.CityID)

	if err != nil {
		return domain.Address{}, err
	}

	return data, nil
}

func (service *AddressesService) Delete(ctx context.Context, userID primitive.ObjectID, addressID primitive.ObjectID) error {
	return service.repo.Delete(ctx, userID, addressID)
}

func NewAddressesService(repo repository.Addresses, areaService Areas) *AddressesService {
	return &AddressesService{
		repo: repo,
		area: areaService,
	}
}
