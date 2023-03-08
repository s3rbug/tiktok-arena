package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"tiktok-arena/configuration"
	"tiktok-arena/database"
	"tiktok-arena/models"
	"time"
)

func RegisterUser(c *fiber.Ctx) error {
	var payload *models.RegisterInput

	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	//	TODO: useless?
	err = models.ValidateStruct(payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	if database.CheckIfUserExists(payload.Name) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User " + payload.Name + " already exists"})
	}

	newUser := models.User{
		Name:     payload.Name,
		Password: string(hashedPassword),
	}

	err = database.CreateNewUser(&newUser)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened"})
	}

	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"data": fiber.Map{
				"user": newUser,
			},
		},
	)
}

func LoginUser(c *fiber.Ctx) error {
	var payload *models.LoginInput

	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	err = models.ValidateStruct(payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	user, err := database.FindUserByName(payload.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	now := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"name": user.Name,
		"exp":  now.Add(configuration.EnvConfig.JwtExpiresIn).Unix(),
		"iat":  now.Unix(),
		"nbf":  now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(configuration.EnvConfig.JwtSecret))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": fmt.Sprintf("generating JWT Token failed: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": tokenString})
}
