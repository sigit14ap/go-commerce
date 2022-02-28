package app

import (
	"context"
	"fmt"
	"github.com/sigit14ap/go-commerce/internal/database/seeds"
	"github.com/sigit14ap/go-commerce/internal/delivery/http/middleware"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/pkg/courier"
	"github.com/sigit14ap/go-commerce/pkg/storage"
	"net/http"
	"time"

	"github.com/sigit14ap/go-commerce/internal/config"
	delivery "github.com/sigit14ap/go-commerce/internal/delivery/http"
	"github.com/sigit14ap/go-commerce/internal/repository"
	"github.com/sigit14ap/go-commerce/internal/service"
	"github.com/sigit14ap/go-commerce/pkg/auth"
	"github.com/sigit14ap/go-commerce/pkg/database/mongodb"
	"github.com/sigit14ap/go-commerce/pkg/database/redis"
	_ "github.com/sigit14ap/go-commerce/pkg/logging"
	log "github.com/sirupsen/logrus"
)

func Run(configPath string, command domain.Command) {
	log.Info("Application start ...")
	log.Info("Logger initialized ...")

	cfg := config.GetConfig(configPath)
	log.Info("Config created ...")

	mongoClient, err := mongodb.NewClient(context.Background(), cfg)

	if err != nil {
		log.Fatal(err)
	}

	log.Info("Mongodb connected ...")
	db := mongoClient.Database(cfg.DB.Database)

	redisClient, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Redis connected ...")

	tokenProvider := auth.NewTokenProvider(cfg, redisClient)
	log.Info("Token provider initialized")

	repos := repository.NewRepositories(db)
	services := service.NewServices(service.Deps{
		Repos:       repos,
		RedisClient: redisClient,
	})

	storageProvider := storage.NewStorageProvider(cfg)

	courierProvider := courier.NewCourierProvider()

	middlewares := middleware.NewMiddleware(services)

	handlers := delivery.NewHandler(services, tokenProvider, storageProvider, courierProvider, middlewares)
	log.Info("Services, repositories and handlers initialized")

	seeder := seeds.NewDatabase(services, courierProvider)
	if command.Seeds {
		seeder.Run()
	}

	server := &http.Server{
		Handler:      handlers.Init(),
		Addr:         fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("Server started on  %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)

	log.Fatal(server.ListenAndServe())
}
