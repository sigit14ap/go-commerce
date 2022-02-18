package service

import (
	"context"
	"fmt"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"github.com/sigit14ap/go-commerce/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	repo        repository.Users
	cartService Carts
}

func NewUsersService(repo repository.Users, cartService Carts) *UsersService {
	return &UsersService{
		repo:        repo,
		cartService: cartService,
	}
}

func (u *UsersService) FindAll(ctx context.Context) ([]domain.User, error) {
	return u.repo.FindAll(ctx)
}

func (u *UsersService) FindByID(ctx context.Context, userID primitive.ObjectID) (domain.User, error) {
	return u.repo.FindByID(ctx, userID)
}

func (u *UsersService) FindByCredentials(ctx context.Context, signInDTO dto.SignInDTO) (domain.LoginUser, error) {
	return u.repo.FindByCredentials(ctx, signInDTO.Email, signInDTO.Password)
}

func (u *UsersService) FindUserInfo(ctx context.Context, userID primitive.ObjectID) (domain.UserInfo, error) {
	return u.repo.FindUserInfo(ctx, userID)
}

func (u UsersService) Create(ctx context.Context, userDTO dto.CreateUserDTO) (domain.User, error) {
	hashPassword, err := HashPassword(userDTO.Password)

	if err != nil {
		return domain.User{}, err
	}

	return u.repo.Create(ctx, domain.User{
		Name:     userDTO.Name,
		Email:    userDTO.Email,
		Password: hashPassword,
	})
}

func (u *UsersService) Update(ctx context.Context, userDTO dto.UpdateUserDTO, userID primitive.ObjectID) (domain.User, error) {
	return u.repo.Update(ctx, dto.UpdateUserInput{
		Name: userDTO.Name,
	}, userID)
}

func (u *UsersService) Delete(ctx context.Context, userID primitive.ObjectID) error {
	_, err := u.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("not found user with id: %s", userID)
	}
	err = u.repo.Delete(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *UsersService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
