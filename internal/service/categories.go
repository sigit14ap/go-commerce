package service

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"github.com/sigit14ap/go-commerce/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoriesService struct {
	repo repository.Categories
}

func (service *CategoriesService) FindAll(ctx context.Context) ([]domain.Category, error) {
	categories, err := service.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (service *CategoriesService) FindByID(ctx context.Context, categoryID primitive.ObjectID) (domain.Category, error) {
	category, err := service.repo.FindByID(ctx, categoryID)
	if err != nil {
		return domain.Category{}, err
	}

	return category, err
}

func (service *CategoriesService) Create(ctx context.Context, category dto.CreateCategoryDTO) (domain.Category, error) {
	return service.repo.Create(ctx, domain.Category{
		Name:        category.Name,
		Description: category.Description,
		Icon:        category.Icon,
	})
}

func (service *CategoriesService) Update(ctx context.Context, categoryDTO dto.UpdateCategoryDTO, categoryID primitive.ObjectID) (domain.Category, error) {
	return service.repo.Update(ctx, dto.UpdateCategoryInput{
		Name:        categoryDTO.Name,
		Description: categoryDTO.Description,
		Icon:        categoryDTO.Icon,
	}, categoryID)
}

func (service *CategoriesService) Delete(ctx context.Context, categoryID primitive.ObjectID) error {
	err := service.repo.Delete(ctx, categoryID)
	return err
}

func NewCategoriesService(repo repository.Categories) *CategoriesService {
	return &CategoriesService{
		repo: repo,
	}
}
