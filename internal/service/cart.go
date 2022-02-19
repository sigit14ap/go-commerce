package service

import (
	"context"
	"fmt"

	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"github.com/sigit14ap/go-commerce/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartService struct {
	repo           repository.Carts
	productService Products
}

func (c *CartService) FindAll(ctx context.Context) ([]domain.Cart, error) {
	carts, err := c.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for i, cart := range carts {
		var totalPrice float64
		for _, cartItem := range cart.CartItems {
			product, err := c.productService.FindByID(ctx, cartItem.ProductID)

			if err != nil {
				return nil, fmt.Errorf("Product no longer exists in stock")
			}

			totalPrice += product.Price * float64(cartItem.Quantity)
		}

		carts[i].TotalPrice = totalPrice
	}

	return carts, nil
}

func (c *CartService) FindByID(ctx context.Context, userID primitive.ObjectID) (domain.Cart, error) {
	cart, err := c.repo.FindByID(ctx, userID)
	if err != nil {
		return domain.Cart{}, err
	}

	var totalPrice float64
	for _, cartItem := range cart.CartItems {
		product, err := c.productService.FindByID(ctx, cartItem.ProductID)

		if err != nil {
			return domain.Cart{}, fmt.Errorf("Product no longer in stock")
		}

		totalPrice += product.Price * float64(cartItem.Quantity)
	}

	cart.TotalPrice = totalPrice

	return cart, nil
}

func (c *CartService) FindCartItems(ctx context.Context, userID primitive.ObjectID) ([]domain.CartItem, error) {
	cartItems, err := c.repo.FindCartItems(ctx, userID)

	if err != nil {
		return []domain.CartItem{}, err
	}

	var itemList []domain.CartItem

	for _, item := range cartItems {
		product, err := c.productService.FindByID(ctx, item.ProductID)

		if err != nil {
			return []domain.CartItem{}, err
		}

		itemList = append(itemList, domain.CartItem{
			Quantity:  item.Quantity,
			ProductID: item.ProductID,
			Product:   product,
		})
	}

	return itemList, nil
}

func (c *CartService) AddCartItem(ctx context.Context, cartItem domain.CartItem, userID primitive.ObjectID) (domain.CartItem, error) {

	product, err := c.productService.FindByID(ctx, cartItem.ProductID)

	if err != nil {
		return domain.CartItem{}, err
	}

	cartItem.Product = product

	_, err = c.repo.AddCartItem(ctx, cartItem, userID)
	return cartItem, err
}

func (c *CartService) UpdateCartItem(ctx context.Context, cartItem domain.CartItem, userID primitive.ObjectID) (domain.CartItem, error) {
	product, err := c.productService.FindByID(ctx, cartItem.ProductID)

	if err != nil {
		return domain.CartItem{}, err
	}

	cartItem.Product = product

	_, err = c.repo.UpdateCartItem(ctx, cartItem, userID)
	return cartItem, err
}

func (c *CartService) DeleteCartItem(ctx context.Context, productID primitive.ObjectID, userID primitive.ObjectID) error {
	return c.repo.DeleteCartItem(ctx, productID, userID)
}

func (c *CartService) ClearCart(ctx context.Context, userID primitive.ObjectID) error {
	return c.repo.ClearCart(ctx, userID)
}

func (c *CartService) Create(ctx context.Context, cartDTO dto.CreateCartDTO) (domain.Cart, error) {
	return c.repo.Create(ctx, domain.Cart{
		CartItems: cartDTO.CartItems,
	})
}

func (c *CartService) Update(ctx context.Context, cartDTO dto.UpdateCartDTO, cartID primitive.ObjectID) (domain.Cart, error) {
	return c.repo.Update(ctx, dto.UpdateCartInput{
		CartItems: cartDTO.CartItems,
	}, cartID)
}

func (c *CartService) Delete(ctx context.Context, cartID primitive.ObjectID) error {
	return c.repo.Delete(ctx, cartID)
}

func NewCartsService(repo repository.Carts, productsService Products) *CartService {
	return &CartService{
		repo:           repo,
		productService: productsService,
	}
}
