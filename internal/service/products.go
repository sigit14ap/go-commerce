package service

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"github.com/sigit14ap/go-commerce/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductsService struct {
	repo              repository.Products
	reviewsService    Reviews
	categoriesService Categories
}

func (p *ProductsService) FindAll(ctx context.Context) ([]domain.Product, error) {
	products, err := p.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for i, product := range products {
		products[i].TotalRating, err = p.reviewsService.GetTotalReviewRating(ctx, product.ID)
		if err != nil {
			return nil, err
		}

		categoryID, err := primitive.ObjectIDFromHex(product.CategoryID)
		if err != nil {
			return nil, err
		}

		products[i].Category, err = p.categoriesService.FindByID(ctx, categoryID)
		if err != nil {
			return nil, err
		}
	}

	return products, nil
}

func (p *ProductsService) FindByID(ctx context.Context, productID primitive.ObjectID) (domain.Product, error) {
	product, err := p.repo.FindByID(ctx, productID)
	if err != nil {
		return domain.Product{}, err
	}

	product.TotalRating, err = p.reviewsService.GetTotalReviewRating(ctx, productID)

	categoryID, err := primitive.ObjectIDFromHex(product.CategoryID)
	if err != nil {
		return domain.Product{}, err
	}

	product.Category, err = p.categoriesService.FindByID(ctx, categoryID)
	if err != nil {
		return domain.Product{}, err
	}

	return product, err
}

func (p *ProductsService) Create(ctx context.Context, product dto.CreateProductDTO) (domain.Product, error) {

	var images []domain.ProductImage

	for _, image := range product.Images {
		images = append(images, domain.ProductImage{
			Image: image,
		})
	}

	return p.repo.Create(ctx, domain.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
		Images:      images,
		Weight:      product.Weight,
	})
}

func (p *ProductsService) Update(ctx context.Context, productDTO dto.UpdateProductDTO, productID primitive.ObjectID) (domain.Product, error) {
	return p.repo.Update(ctx, dto.UpdateProductInput{
		Name:        productDTO.Name,
		Description: productDTO.Description,
		Price:       productDTO.Price,
		CategoryID:  productDTO.CategoryID,
	}, productID)
}

func (p *ProductsService) Delete(ctx context.Context, productID primitive.ObjectID) error {
	err := p.repo.Delete(ctx, productID)
	if err != nil {
		return err
	}
	return p.reviewsService.DeleteByProductID(ctx, productID)
}

func NewProductsService(repo repository.Products, reviewsService Reviews, categoriesService Categories) *ProductsService {
	return &ProductsService{
		repo:              repo,
		reviewsService:    reviewsService,
		categoriesService: categoriesService,
	}
}
