package main

import (
	"fmt"
	"log"

	"github.com/fazriegi/fintrack-be/internal/infrastructure/config"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/database"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/logger"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/router.go"
	"github.com/fazriegi/fintrack-be/internal/middleware"
	"github.com/fazriegi/fintrack-be/internal/pkg"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	viperConfig := config.NewViper()
	database.NewMysql(viperConfig)
	jwt := pkg.InitJWT(config.GetString("jwt.key"), config.GetUint("jwt.accessToken.expMinute"), config.GetUint("jwt.refreshToken.expDay"))
	file := logger.New(viperConfig)
	defer file.Close()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(middleware.LogMiddleware())
	port := viperConfig.GetInt("web.port")
	router.NewRoute(app, jwt)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
