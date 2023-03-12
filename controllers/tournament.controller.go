package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"math"
	"tiktok-arena/database"
	"tiktok-arena/models"
)

// @Summary		Create new tournament
// @Description	Create new tournament for current user
// @Tags			tournament
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			payload	body		models.CreateTournament	true	"Data to create tournament"
// @Success		200		{object}	MessageResponseType		"Tournament created"
// @Failure		400		{object}	MessageResponseType		"Error during tournament creation"
// @Router			/tournament [post]
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

// @Summary		Tournament details
// @Description	Get tournament details by its id
// @Tags			tournament
// @Accept			json
// @Produce		json
// @Param			tournamentId	path		string				true	"Tournament id"
// @Success		200				{object}	models.Tournament	"Tournament"
// @Failure		400				{object}	MessageResponseType	"Tournament not found"
// @Router			/tournament/{tournamentId} [get]
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

// @Summary		Tournament tiktoks
// @Description	Get tournament tiktoks
// @Tags			tournament
// @Accept			json
// @Produce		json
// @Param			tournamentId	path		string				true	"Tournament id"
// @Success		200				{array}		models.Tiktok		"Tournament tiktoks"
// @Failure		400				{object}	MessageResponseType	"Tournament not found"
// @Router			/tournament/{tournamentId}/tiktoks [get]
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

// @Summary		Tournament contest
// @Description	Get tournament contest
// @Tags			tournament
// @Accept			json
// @Produce		json
// @Param			tournamentId	path		string					true	"Tournament id"
// @Param			payload			body		models.ContestPayload	true	"Contest type"
// @Success		200				{array}		[]models.ContestItem	"Contest items, each array represents round of contest"
// @Failure		400				{object}	MessageResponseType		"Failed to return tournament contest"
// @Router			/tournament/{tournamentId}/contest [get]
func GetTournamentContest(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")
	if tournamentId == "" {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not a valid tournament id", tournamentId))
	}
	var payload *models.ContestPayload

	err := c.BodyParser(&payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = models.ValidateStruct(payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if !models.CheckIfAllowedTournamentType(payload.ContestType) {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not allowed tournament format", payload.ContestType),
		)
	}
	tiktoks, err := database.GetTournamentTiktoksById(tournamentId)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Could not get tiktoks for tournament with id %s", tournamentId))
	}

	if payload.ContestType == "single elimination" {
		return c.Status(fiber.StatusOK).JSON(SingleElimination(tiktoks))
	}
	return MessageResponse(c, fiber.StatusBadRequest, "Unknown error")
}

func SingleElimination(tiktoks []models.Tiktok) models.SingleEliminationBracket {
	countTiktok := len(tiktoks)
	countRound := math.Ceil(math.Log2(float64(countTiktok)))
	countFirstRoundTiktoks := (countTiktok - int(math.Pow(2, countRound-1))) * 2

	var bracket models.SingleEliminationBracket
	var rounds []models.Round

	var firstRoundMatches []models.Match
	var secondRoundMatches []models.Match
	// This slice should store MatchOption or TiktokOption
	var secondRoundParticipators []any

	// Filling first round with firstRoundMatches and appending MatchOptions to second round participators
	for j := 0; j < countFirstRoundTiktoks; j += 2 {
		matchID := uuid.NewString()
		firstRoundMatches = append(firstRoundMatches, models.Match{
			MatchID: matchID,
			FirstOption: models.TiktokOption{
				TiktokURL: tiktoks[j].URL,
			},
			SecondOption: models.TiktokOption{
				TiktokURL: tiktoks[j+1].URL,
			},
		})
		secondRoundParticipators = append(secondRoundParticipators, models.MatchOption{MatchID: matchID})
	}
	// Appending first round firstRoundMatches to rounds
	rounds = append(rounds, models.Round{
		Round:   1,
		Matches: firstRoundMatches,
	})
	// Appending TiktokOptions to second round participators
	for i := 0; i < countTiktok-countFirstRoundTiktoks; i++ {
		secondRoundParticipators = append(secondRoundParticipators, models.TiktokOption{TiktokURL: tiktoks[countTiktok-1-i].URL})
	}
	// Generating second round firstRoundMatches
	for i := 0; i < countTiktok-countFirstRoundTiktoks; i += 2 {
		match := models.Match{
			MatchID:      uuid.NewString(),
			FirstOption:  secondRoundParticipators[i],
			SecondOption: secondRoundParticipators[i+1],
		}
		secondRoundMatches = append(secondRoundMatches, match)
	}
	// Generating second round
	secondRound := models.Round{
		Round:   2,
		Matches: secondRoundMatches}
	rounds = append(rounds, secondRound)

	previousRoundMatches := secondRoundMatches
	for roundID := 3; roundID <= int(countRound); roundID++ {
		// Generating Nth round matches (where N > 2)
		var currentRoundMatches []models.Match
		for matchID := 0; matchID < len(previousRoundMatches); matchID += 2 {
			match := models.Match{
				MatchID:      uuid.NewString(),
				FirstOption:  models.MatchOption{MatchID: previousRoundMatches[matchID].MatchID},
				SecondOption: models.MatchOption{MatchID: previousRoundMatches[matchID+1].MatchID},
			}
			currentRoundMatches = append(currentRoundMatches, match)
		}
		// Generating Nth round (where N > 2)
		round := models.Round{
			Round:   roundID,
			Matches: currentRoundMatches,
		}
		rounds = append(rounds, round)

		previousRoundMatches = currentRoundMatches
	}
	// Creating bracket
	bracket = models.SingleEliminationBracket{
		CountMatches: countTiktok - 1,
		Rounds:       &rounds,
	}
	return bracket
}

func GetAllTournaments(c *fiber.Ctx) error {
	tournaments, err := database.GetAllTournaments()
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, "Failed to get tournaments")
	}
	return c.Status(fiber.StatusOK).JSON(tournaments)
}
