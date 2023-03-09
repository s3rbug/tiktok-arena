package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"tiktok-arena/database"
	"tiktok-arena/models"
	"tiktok-arena/utils"
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
		return utils.FiberMessage(c, fiber.StatusBadRequest, err.Error())
	}

	var payload *models.CreateTournament

	err = c.BodyParser(&payload)
	if err != nil {
		return utils.FiberMessage(c, fiber.StatusBadRequest, err.Error())
	}

	err = models.ValidateStruct(payload)
	if err != nil {
		return utils.FiberMessage(c, fiber.StatusBadRequest, err.Error())
	}

	if !models.CheckIfAllowedTournamentSize(payload.Size) {
		return utils.FiberMessage(c, fiber.StatusBadRequest,
			fmt.Sprintf("%d is incorrect tournament size", payload.Size))
	}

	if database.CheckIfTournamentExists(payload.Name) {
		return utils.FiberMessage(c, fiber.StatusBadRequest,
			fmt.Sprintf("Tournament %s already exists", payload.Name))
	}

	newTournamentId, err := uuid.NewRandom()
	if err != nil {
		return utils.FiberMessage(c, fiber.StatusBadRequest, err.Error())
	}

	newTournament := models.Tournament{
		ID:     &newTournamentId,
		Name:   payload.Name,
		UserID: &userId,
		Size:   payload.Size,
	}
	err = database.CreateNewTournament(&newTournament)
	if err != nil {
		return utils.FiberMessage(c, fiber.StatusBadRequest, err.Error())
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
			return utils.FiberMessage(c, fiber.StatusBadRequest, err.Error())
		}
	}

	return utils.FiberMessage(c, fiber.StatusOK,
		fmt.Sprintf("Successfully created tournament %s", payload.Name))
}

func GetTournamentDetails(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")
	if tournamentId == "" {
		return utils.FiberMessage(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not a valid tournament id", tournamentId))
	}
	tournament, err := database.GetTournamentById(tournamentId)
	if err != nil {
		return utils.FiberMessage(c, fiber.StatusBadRequest,
			fmt.Sprintf("Could not get tournament with id %s", tournamentId))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tournament": tournament,
	})
}

func GetTournamentTiktoks(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")
	if tournamentId == "" {
		return utils.FiberMessage(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not a valid tournament id", tournamentId))
	}
	tiktoks, err := database.GetTournamentTiktoksById(tournamentId)
	if err != nil {
		return utils.FiberMessage(c, fiber.StatusBadRequest,
			fmt.Sprintf("Could not get tiktoks for tournament with id %s", tournamentId))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tiktoks": tiktoks,
	})
}
