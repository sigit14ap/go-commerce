package repository

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"github.com/sigit14ap/go-commerce/pkg/courier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AreasRepo struct {
	dbProvince *mongo.Collection
	dbCity     *mongo.Collection
}

func (area *AreasRepo) FindProvinceByThirdParty(ctx context.Context, provinceID string) (domain.Province, error) {
	result := area.dbProvince.FindOne(ctx, bson.M{"third_party_id": provinceID})

	var province domain.Province
	err := result.Decode(&province)

	return province, err
}

func (area *AreasRepo) GetProvinces(ctx context.Context) ([]domain.Province, error) {
	provinces := []domain.Province{}

	cursor, err := area.dbProvince.Find(ctx, bson.M{})
	if err != nil {
		return provinces, err
	}

	err = cursor.All(ctx, &provinces)

	return provinces, err
}

func (area *AreasRepo) CreateProvinces(ctx context.Context, province domain.Province) (domain.Province, error) {
	_, err := area.dbProvince.InsertOne(ctx, province)
	return province, err
}

func (area *AreasRepo) GetCities(ctx context.Context, cityListDTO dto.CityListDTO) ([]dto.ThirdPartyCityDTO, error) {
	cities, err := courier.GetCities()
	return cities, err
}

func (area *AreasRepo) CreateCity(ctx context.Context, city domain.City) (domain.City, error) {
	_, err := area.dbCity.InsertOne(ctx, city)
	return city, err
}

func NewAreasRepo(db *mongo.Database) *AreasRepo {
	return &AreasRepo{
		dbProvince: db.Collection(provincesCollection),
		dbCity:     db.Collection(citiesCollection),
	}
}
