package controller

import (
	"net/http"

	"github.com/fazriegi/fintrack-be/internal/entity"
	"github.com/fazriegi/fintrack-be/internal/infrastructure/logger"
	"github.com/fazriegi/fintrack-be/internal/pkg"
	"github.com/fazriegi/fintrack-be/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AssetController interface {
	ListAssetCategory(ctx *fiber.Ctx) error
}

type assetController struct {
	usecase usecase.AssetUsecase
	logger  *logrus.Logger
}

func NewAssetController(usecase usecase.AssetUsecase) AssetController {
	logger := logger.Get()
	return &assetController{
		usecase,
		logger,
	}
}

func (c *assetController) ListAssetCategory(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals("user").(entity.User)
	if !ok {
		return ctx.Status(http.StatusUnauthorized).JSON(pkg.NewResponse(http.StatusUnauthorized, pkg.ErrNotAuthorized.Error(), nil, nil))
	}

	response := c.usecase.ListAssetCategory(user.ID)

	return ctx.Status(response.Status.Code).JSON(response)
}
