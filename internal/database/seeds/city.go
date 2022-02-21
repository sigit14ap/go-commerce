package seeds

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/service"
	"github.com/sigit14ap/go-commerce/pkg/courier"
	log "github.com/sirupsen/logrus"
)

type CitySeeder struct {
	services *service.Services
}

func (seeder *CitySeeder) Run(context context.Context) error {
	log.Warnf("City Seeder running ...")

	cities, err := courier.GetCities()

	for _, city := range cities {
		_, err = seeder.services.Areas.CreateCity(context, city)
	}

	return err
}

func NewCitySeeder(services *service.Services) *CitySeeder {
	return &CitySeeder{
		services: services,
	}
}
