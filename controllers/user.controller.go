package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"tiktok-arena/database"
	"tiktok-arena/models"
)

//	@Summary		Create new tournament
//	@Description	Create new tournament for current user
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			payload	body		models.CreateTournament	true	"Data to create tournament"
//	@Success		200		{object}	MessageResponseType		"Tournament created"
//	@Failure		400		{object}	MessageResponseType		"Error during tournament creation"
//	@Router			/tournament [post]
func CreateTournament(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	userId, err := uuid.Parse(claims["sub"].(string))

	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	var payload *models.CreateTournament

	err = c.BodyParser(&payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = models.ValidateStruct(payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if !models.CheckIfAllowedTournamentSize(payload.Size) {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%d is incorrect tournament size", payload.Size))
	}

	if payload.Size != len(payload.Tiktoks) {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Tournament size and count of tiktoks mismatch (%d != %d)",
				payload.Size,
				len(payload.Tiktoks)),
		)
	}

	if database.CheckIfTournamentExists(payload.Name) {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Tournament %s already exists", payload.Name))
	}

	newTournamentId, err := uuid.NewRandom()
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	newTournament := models.Tournament{
		ID:     &newTournamentId,
		Name:   payload.Name,
		UserID: &userId,
		Size:   payload.Size,
	}
	err = database.CreateNewTournament(&newTournament)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
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
			return MessageResponse(c, fiber.StatusBadRequest, err.Error())
		}
	}

	return MessageResponse(c, fiber.StatusOK,
		fmt.Sprintf("Successfully created tournament %s", payload.Name))
}

//	@Summary		Tournament details
//	@Description	Get tournament details by its id
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Param			tournamentId	path		string				true	"Tournament id"
//	@Success		200				{object}	models.Tournament	"Tournament"
//	@Failure		400				{object}	MessageResponseType	"Tournament not found"
//	@Router			/tournament/{tournamentId} [get]
func GetTournamentDetails(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")
	if tournamentId == "" {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not a valid tournament id", tournamentId))
	}
	tournament, err := database.GetTournamentById(tournamentId)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Could not get tournament with id %s", tournamentId))
	}
	return c.Status(fiber.StatusOK).JSON(tournament)
}

//	@Summary		Tournament tiktoks
//	@Description	Get tournament tiktoks
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Param			tournamentId	path		string				true	"Tournament id"
//	@Success		200				{array}		models.Tiktok		"Tournament tiktoks"
//	@Failure		400				{object}	MessageResponseType	"Tournament not found"
//	@Router			/tournament/{tournamentId}/tiktoks [get]
func GetTournamentTiktoks(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")
	if tournamentId == "" {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not a valid tournament id", tournamentId))
	}
	tiktoks, err := database.GetTournamentTiktoksById(tournamentId)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Could not get tiktoks for tournament with id %s", tournamentId))
	}
	return c.Status(fiber.StatusOK).JSON(tiktoks)
}
