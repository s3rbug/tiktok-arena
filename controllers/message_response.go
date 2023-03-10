package controllers

import "github.com/gofiber/fiber/v2"

type MessageResponseType struct {
	Message string
}

func MessageResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
	})
}
