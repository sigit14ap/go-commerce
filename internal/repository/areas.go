package repository

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"github.com/sigit14ap/go-commerce/pkg/courier"
)

type AreasRepo struct {
}

func (area *AreasRepo) GetProvinces(ctx context.Context) ([]domain.Province, error) {
	provinces, err := courier.GetProvinces()
	return provinces, err
}

func (area *AreasRepo) GetCities(ctx context.Context, cityListDTO dto.CityListDTO) ([]domain.City, error) {
	cities, err := courier.GetCities(cityListDTO)
	return cities, err
}

func NewAreasRepo() *AreasRepo {
	return &AreasRepo{}
}
