package main

import (
	"github.com/joho/godotenv"
	"github.com/sigit14ap/go-commerce/internal/app"
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

	//var param1 string
	//flag.StringVar(&param1, "param1", "", "Parameter 1")
	//flag.Parse()
	//fmt.Print("Missing required parameter 1 : %s", param1)

	app.Run("config/config.yml")
}
