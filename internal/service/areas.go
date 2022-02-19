package service

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"github.com/sigit14ap/go-commerce/internal/repository"
)

type AreaService struct {
	repo repository.Areas
}

func (area *AreaService) GetProvinces(ctx context.Context) ([]domain.Province, error) {
	return area.repo.GetProvinces(ctx)
}

func (area *AreaService) GetCities(ctx context.Context, cityListDTO dto.CityListDTO) ([]domain.City, error) {
	return area.repo.GetCities(ctx, cityListDTO)
}

func NewAreasService(repo repository.Areas) *AreaService {
	return &AreaService{
		repo: repo,
	}
}
