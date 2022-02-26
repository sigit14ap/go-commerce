package seeds

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/service"
	"github.com/sigit14ap/go-commerce/pkg/courier"
	log "github.com/sirupsen/logrus"
)

type CitySeeder struct {
	services *service.Services
	courier  *courier.Provider
}

func (seeder *CitySeeder) Run(context context.Context) error {
	log.Warnf("City Seeder running ...")

	cities, err := seeder.courier.GetCities()
	for _, city := range cities {

		province, err := seeder.services.Areas.FindProvinceByThirdParty(context, city.ProvinceID)

		if err != nil {
			return err
		}

		dataCity := domain.City{
			ProvinceID:   province.ID,
			ThirdPartyID: city.CityID,
			Name:         city.Name,
		}

		_, err = seeder.services.Areas.CreateCity(context, dataCity)
	}

	return err
}

func NewCitySeeder(services *service.Services, courier *courier.Provider) *CitySeeder {
	return &CitySeeder{
		services: services,
		courier:  courier,
	}
}
