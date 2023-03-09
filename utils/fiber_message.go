package utils

import "github.com/gofiber/fiber/v2"

func FiberMessage(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
	})
}
