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
	ListAsset(ctx *fiber.Ctx) error
	SubmitAsset(ctx *fiber.Ctx) error
	GetById(ctx *fiber.Ctx) error
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

func (c *assetController) ListAsset(ctx *fiber.Ctx) error {
	var reqQuery entity.ListAssetRequest

	if err := ctx.QueryParser(&reqQuery); err != nil {
		c.logger.Errorf("error parsing query param: %s", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(pkg.NewResponse(http.StatusBadRequest, pkg.ErrParseQueryParam.Error(), nil, nil))
	}

	user, ok := ctx.Locals("user").(entity.User)
	if !ok {
		return ctx.Status(http.StatusUnauthorized).JSON(pkg.NewResponse(http.StatusUnauthorized, pkg.ErrNotAuthorized.Error(), nil, nil))
	}

	reqQuery.UserId = user.ID

	response := c.usecase.ListAsset(reqQuery)

	return ctx.Status(response.Status.Code).JSON(response)
}

func (c *assetController) SubmitAsset(ctx *fiber.Ctx) error {
	var reqBody entity.SubmitAssetRequest

	if err := ctx.BodyParser(&reqBody); err != nil {
		c.logger.Errorf("error parsing request body: %s", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(pkg.NewResponse(http.StatusBadRequest, pkg.ErrParseQueryParam.Error(), nil, nil))
	}

	// validate reqBody struct
	validationErr := pkg.ValidateRequest(&reqBody)
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		return ctx.Status(http.StatusUnprocessableEntity).JSON(pkg.NewResponse(http.StatusUnprocessableEntity, pkg.ErrValidation.Error(), errResponse, nil))
	}

	user, ok := ctx.Locals("user").(entity.User)
	if !ok {
		return ctx.Status(http.StatusUnauthorized).JSON(pkg.NewResponse(http.StatusUnauthorized, pkg.ErrNotAuthorized.Error(), nil, nil))
	}

	reqBody.UserId = user.ID

	response := c.usecase.SubmitAsset(reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}

func (c *assetController) GetById(ctx *fiber.Ctx) error {
	var reqQuery entity.GetAssetByIdRequest

	if err := ctx.ParamsParser(&reqQuery); err != nil {
		c.logger.Errorf("error parsing query param: %s", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(pkg.NewResponse(http.StatusBadRequest, pkg.ErrParseQueryParam.Error(), nil, nil))
	}

	user, ok := ctx.Locals("user").(entity.User)
	if !ok {
		return ctx.Status(http.StatusUnauthorized).JSON(pkg.NewResponse(http.StatusUnauthorized, pkg.ErrNotAuthorized.Error(), nil, nil))
	}

	reqQuery.UserId = user.ID

	response := c.usecase.GetById(reqQuery)

	return ctx.Status(response.Status.Code).JSON(response)
}
