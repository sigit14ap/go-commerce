package service

import (
	"context"

	"github.com/go-redis/redis/v7"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"github.com/sigit14ap/go-commerce/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users interface {
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByID(ctx context.Context, userID primitive.ObjectID) (domain.User, error)
	FindByCredentials(ctx context.Context, signInDTO dto.SignInDTO) (domain.LoginUser, error)
	FindUserInfo(ctx context.Context, userID primitive.ObjectID) (domain.UserInfo, error)
	Create(ctx context.Context, userDTO dto.CreateUserDTO) (domain.User, error)
	Update(ctx context.Context, userDTO dto.UpdateUserDTO,
		userID primitive.ObjectID) (domain.User, error)
	Delete(ctx context.Context, userID primitive.ObjectID) error
	CheckPasswordHash(password, hash string) bool
}

type Products interface {
	FindAll(ctx context.Context) ([]domain.Product, error)
	FindByID(ctx context.Context, productID primitive.ObjectID) (domain.Product, error)
	Create(ctx context.Context, productDTO dto.CreateProductDTO) (domain.Product, error)
	Update(ctx context.Context, productDTO dto.UpdateProductDTO,
		productID primitive.ObjectID) (domain.Product, error)
	Delete(ctx context.Context, productID primitive.ObjectID) error
}

type Reviews interface {
	FindAll(ctx context.Context) ([]domain.Review, error)
	FindByID(ctx context.Context, reviewID primitive.ObjectID) (domain.Review, error)
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Review, error)
	FindByProductID(ctx context.Context, productID primitive.ObjectID) ([]domain.Review, error)
	GetTotalReviewRating(ctx context.Context, productID primitive.ObjectID) (float64, error)
	Create(ctx context.Context, review dto.CreateReviewInput) (domain.Review, error)
	Delete(ctx context.Context, reviewID primitive.ObjectID) error
	DeleteByProductID(ctx context.Context, productID primitive.ObjectID) error
}

type Admins interface {
	FindByCredentials(ctx context.Context, signInDTO dto.SignInDTO) (domain.Admin, error)
	CheckPasswordHash(password, hash string) bool
}

type Carts interface {
	FindAll(ctx context.Context) ([]domain.Cart, error)
	FindByID(ctx context.Context, userID primitive.ObjectID) (domain.Cart, error)
	FindCartItems(ctx context.Context, userID primitive.ObjectID) ([]domain.CartItem, error)
	AddCartItem(ctx context.Context, cartItem domain.CartItem, userID primitive.ObjectID) (domain.CartItem, error)
	UpdateCartItem(ctx context.Context, cartItem domain.CartItem, userID primitive.ObjectID) (domain.CartItem, error)
	DeleteCartItem(ctx context.Context, productID primitive.ObjectID, userID primitive.ObjectID) error
	ClearCart(ctx context.Context, userID primitive.ObjectID) error
	Create(ctx context.Context, cartDTO dto.CreateCartDTO) (domain.Cart, error)
	Update(ctx context.Context, cartDTO dto.UpdateCartDTO,
		cartID primitive.ObjectID) (domain.Cart, error)
	Delete(ctx context.Context, cartID primitive.ObjectID) error
}

type Orders interface {
	FindAll(ctx context.Context) ([]domain.Order, error)
	FindByID(ctx context.Context, orderID primitive.ObjectID) (domain.Order, error)
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Order, error)
	Create(ctx context.Context, orderDTO dto.CreateOrderDTO) (domain.Order, error)
	Update(ctx context.Context, orderDTO dto.UpdateOrderDTO,
		orderID primitive.ObjectID) (domain.Order, error)
	Delete(ctx context.Context, orderID primitive.ObjectID) error
}

type Payment interface {
	GetPaymentLink(ctx context.Context, orderID primitive.ObjectID) (string, error)
}

type Categories interface {
	FindAll(ctx context.Context) ([]domain.Category, error)
	FindByID(ctx context.Context, categoryID primitive.ObjectID) (domain.Category, error)
	Create(ctx context.Context, categoryDTO dto.CreateCategoryDTO) (domain.Category, error)
	Update(ctx context.Context, categoryDTO dto.UpdateCategoryDTO,
		categoryID primitive.ObjectID) (domain.Category, error)
	Delete(ctx context.Context, categoryID primitive.ObjectID) error
}

type Areas interface {
	GetProvinces(ctx context.Context) ([]domain.Province, error)
	GetCities(ctx context.Context, cityListDTO dto.CityListDTO) ([]domain.City, error)
}

type Addresses interface {
	FindAll(ctx context.Context, userID primitive.ObjectID) ([]domain.Address, error)
}

type Services struct {
	Users      Users
	Products   Products
	Reviews    Reviews
	Admins     Admins
	Carts      Carts
	Orders     Orders
	Payment    Payment
	Categories Categories
	Areas      Areas
	Addresses  Addresses
}

type Deps struct {
	Repos       *repository.Repositories
	Services    *Services
	RedisClient *redis.Client
}

func NewServices(deps Deps) *Services {
	reviewsService := NewReviewsService(deps.Repos.Reviews, deps.RedisClient)
	CategoriesService := NewCategoriesService(deps.Repos.Categories)
	productsService := NewProductsService(deps.Repos.Products, reviewsService, CategoriesService)
	adminsService := NewAdminsService(deps.Repos.Admins)
	cartsService := NewCartsService(deps.Repos.Carts, productsService)
	usersService := NewUsersService(deps.Repos.Users, cartsService)
	ordersService := NewOrdersService(deps.Repos.Orders, productsService, cartsService)
	areaService := NewAreasService(deps.Repos.Areas)
	addressService := NewAddressesService(deps.Repos.Addresses)
	// paymentService := NewPaymentService(ordersService, productsService)

	return &Services{
		Users:      usersService,
		Products:   productsService,
		Reviews:    reviewsService,
		Admins:     adminsService,
		Carts:      cartsService,
		Orders:     ordersService,
		Categories: CategoriesService,
		Areas:      areaService,
		Addresses:  addressService,
		// Payment:  paymentService,
	}
}
