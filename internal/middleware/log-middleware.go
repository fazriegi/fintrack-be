package middleware

import (
	"github.com/fazriegi/fintrack-be/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v2"
)

func LogMiddleware() func(ctx *fiber.Ctx) error {
	log := logger.Get()

	return func(ctx *fiber.Ctx) error {
		if err := ctx.Next(); err != nil {
			log.Errorf("error handling request: %s | method=%s | uri=%s", err.Error(), ctx.Method(), ctx.OriginalURL())
			return err
		}

		log.Infof("incoming request | method=%s | uri=%s", ctx.Method(), ctx.OriginalURL())
		return nil
	}
}
