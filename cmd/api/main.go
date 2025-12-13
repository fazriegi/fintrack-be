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

	origins := config.GetString("cors.origins")

	app.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Use(middleware.LogMiddleware())
	port := config.GetInt("web.port")
	router.NewRoute(app, jwt)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
