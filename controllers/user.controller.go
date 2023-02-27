package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func WhoAmI(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":       claims["sub"],
		"username": claims["name"],
	})
}
