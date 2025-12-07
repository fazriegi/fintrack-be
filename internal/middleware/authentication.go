package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fazriegi/fintrack-be/internal/entity"
	"github.com/fazriegi/fintrack-be/internal/pkg"
	"github.com/gofiber/fiber/v2"
)

func Authentication(jwt *pkg.JWT) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var response = pkg.Response{}
		header := ctx.Get("Authorization")
		isHasBearer := strings.HasPrefix(header, "Bearer")

		if !isHasBearer {
			response = pkg.NewResponse(http.StatusUnauthorized, "sign in to proceed", nil, nil)

			return ctx.Status(response.Code).JSON(response)
		}

		tokenString := strings.Split(header, " ")[1]

		verifiedToken, err := jwt.VerifyJWTTOken(tokenString)
		if err != nil {
			response = pkg.NewResponse(http.StatusUnauthorized, err.Error(), nil, nil)
			return ctx.Status(response.Code).JSON(response)
		}

		jsonData, _ := json.Marshal(verifiedToken)

		var user entity.User
		json.Unmarshal(jsonData, &user)

		ctx.Locals("user", user)

		return ctx.Next()
	}
}
