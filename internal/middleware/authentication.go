package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/fazriegi/fintrack-be/internal/entity"
	"github.com/fazriegi/fintrack-be/internal/pkg"
	"github.com/gofiber/fiber/v2"
)

func Authentication(jwt *pkg.JWT) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var response = pkg.Response{}

		tokenString := ctx.Cookies("access_token")

		verifiedToken, err := jwt.VerifyToken(tokenString, "access")
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
