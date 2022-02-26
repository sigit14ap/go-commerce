package seeds

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/service"
	"github.com/sigit14ap/go-commerce/pkg/courier"
)

type Function interface {
	Run(context context.Context) error
}

type DatabaseSeeder struct {
	Province Function
	City     Function
}

func (seeds *DatabaseSeeder) Run() {
	var context context.Context
	_ = seeds.Province.Run(context)
	_ = seeds.City.Run(context)
}

func NewDatabase(services *service.Services, courier *courier.Provider) *DatabaseSeeder {
	return &DatabaseSeeder{
		Province: NewProvinceSeeder(services, courier),
		City:     NewCitySeeder(services, courier),
	}
}
