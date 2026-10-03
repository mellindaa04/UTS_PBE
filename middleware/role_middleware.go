package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")

		if userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "akses ditolak",
			})
		}

		return c.Next()
	}
}
