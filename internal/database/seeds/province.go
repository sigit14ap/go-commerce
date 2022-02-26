package seeds

import (
	"context"
	"github.com/sigit14ap/go-commerce/internal/service"
	"github.com/sigit14ap/go-commerce/pkg/courier"
	log "github.com/sirupsen/logrus"
)

type ProvinceSeeder struct {
	services *service.Services
	courier  *courier.Provider
}

func (seeder *ProvinceSeeder) Run(context context.Context) error {
	log.Warnf("Province Seeder running ...")

	provinces, err := seeder.courier.GetProvinces()

	for _, province := range provinces {
		_, err = seeder.services.Areas.CreateProvinces(context, province)
	}

	return err
}

func NewProvinceSeeder(services *service.Services, courier *courier.Provider) *ProvinceSeeder {
	return &ProvinceSeeder{
		services: services,
		courier:  courier,
	}
}
