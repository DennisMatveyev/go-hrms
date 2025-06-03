package admin

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return fiber.ErrUnauthorized
		}
		if strings.HasPrefix(auth, "Token ") {
			password := strings.TrimPrefix(auth, "Token ")
			if password == os.Getenv("ADMIN_PASSWORD") {
				return c.Next()
			}
			return fiber.ErrUnauthorized
		}
		return fiber.ErrUnauthorized
	}
}
