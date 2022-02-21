package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/sigit14ap/go-commerce/internal/app"
	"github.com/sigit14ap/go-commerce/internal/domain"
	log "github.com/sirupsen/logrus"
)

// @title        E-commerce API
// @version      1.0
// @description  This is simple api of e-commerce shop

// @contact.name   API Support
// @contact.url    https://t.me/paw1a
// @contact.email  paw1a@yandex.ru

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  AdminAuth
// @in                          header
// @name                        Authorization

// @securityDefinitions.apikey  UserAuth
// @in                          header
// @name                        Authorization
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var seeds bool

	// flags declaration using flag package
	flag.BoolVar(&seeds, "seeds", false, "Running seeders")

	flag.Parse() // after declaring flags we need to call it

	command := domain.Command{
		Seeds: seeds,
	}

	app.Run("config/config.yml", command)
}
