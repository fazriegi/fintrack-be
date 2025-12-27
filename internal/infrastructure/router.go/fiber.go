package router

import (
	"github.com/fazriegi/fintrack-be/internal/controller"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/repository"
	"github.com/fazriegi/fintrack-be/internal/middleware"
	"github.com/fazriegi/fintrack-be/internal/pkg"
	"github.com/fazriegi/fintrack-be/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App, jwt *pkg.JWT) {
	userRepo := repository.NewUserRepository()
	authRepo := repository.NewAuthRepository()
	authUC := usecase.NewAuthUsecase(userRepo, authRepo, jwt)
	authController := controller.NewAuthController(authUC)

	v1 := app.Group("/api/v1")
	{
		v1.Post("/register", authController.Register)
		v1.Post("/login", authController.Login)
		v1.Get("/check-token", middleware.Authentication(jwt), authController.CheckToken)
		v1.Post("/refresh-token", authController.RefreshToken)
		v1.Post("/logout", authController.Logout)
	}

	assetRepo := repository.NewAssetRepository()
	assetUC := usecase.NewAssetUsecase(assetRepo, jwt)
	assetController := controller.NewAssetController(assetUC)
	asset := v1.Group("asset")
	{
		asset.Get("/list-category", middleware.Authentication(jwt), assetController.ListAssetCategory)
		asset.Post("/submit", middleware.Authentication(jwt), assetController.SubmitAsset)
		asset.Get("/list", middleware.Authentication(jwt), assetController.ListAsset)
		asset.Get("/:id", middleware.Authentication(jwt), assetController.GetById)
	}
}
