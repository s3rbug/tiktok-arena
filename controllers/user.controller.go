package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"strconv"
	"tiktok-arena/database"
	"tiktok-arena/models"
)

func WhoAmI(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":       claims["sub"],
		"username": claims["name"],
	})
}

func CreateTournament(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	userId, err := uuid.Parse(claims["sub"].(string))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	var payload *models.CreateTournament

	err = c.BodyParser(&payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	//	TODO: useless?
	err = models.ValidateStruct(payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	if !models.CheckIfAllowedTournamentSize(payload.Size) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": strconv.Itoa(payload.Size) + " is incorrect tournament size",
		})
	}

	if database.CheckIfTournamentExists(payload.Name) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Tournament " + payload.Name + " already exists"})
	}

	newTournamentId, err := uuid.NewRandom()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	newTournament := models.Tournament{
		ID:     &newTournamentId,
		Name:   payload.Name,
		UserID: &userId,
		Size:   payload.Size,
	}
	err = database.CreateNewTournament(&newTournament)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	for _, value := range payload.Tiktoks {
		tiktok := models.Tiktok{
			TournamentID: &newTournamentId,
			URL:          value.URL,
			Wins:         0,
			AvgPoints:    0,
		}
		err = database.CreateNewTiktok(&tiktok)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully created tournament " + payload.Name,
	})
}
