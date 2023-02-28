package jwt_middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/urosradivojevic/health/pkg/message"
)

func AuthRequired() fiber.Handler {
	return jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(message.Msg{
				Message: "Unauthorized",
			})
		},
		SigningKey: []byte(os.Getenv("JWT_KEY")),
	})
}
