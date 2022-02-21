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

func (area *AreaService) FindProvinceByThirdParty(ctx context.Context, provinceID string) (domain.Province, error) {
	return area.repo.FindProvinceByThirdParty(ctx, provinceID)
}

func (area *AreaService) GetProvinces(ctx context.Context) ([]domain.Province, error) {
	return area.repo.GetProvinces(ctx)
}

func (area *AreaService) CreateProvinces(ctx context.Context, province domain.Province) (domain.Province, error) {
	return area.repo.CreateProvinces(ctx, province)
}

func (area *AreaService) GetCities(ctx context.Context, cityListDTO dto.CityListDTO) ([]dto.ThirdPartyCityDTO, error) {
	return area.repo.GetCities(ctx, cityListDTO)
}

func (area *AreaService) CreateCity(ctx context.Context, city domain.City) (domain.City, error) {
	return area.repo.CreateCity(ctx, city)
}

func NewAreasService(repo repository.Areas) *AreaService {
	return &AreaService{
		repo: repo,
	}
}
