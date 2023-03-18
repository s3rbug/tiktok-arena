package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"math"
	"math/rand"
	"tiktok-arena/database"
	"tiktok-arena/models"
	"time"
)

// CreateTournament
//
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

// GetTournamentDetails
//
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

// GetTournamentTiktoks
//
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

// GetTournamentContest
//
//	@Summary		Tournament contest
//	@Description	Get tournament contest
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Param			tournamentId	path		string					true	"Tournament id"
//	@Param			payload			query		models.ContestPayload	true	"Contest type"
//	@Success		200				{object}	models.Bracket			"Contest bracket"
//	@Failure		400				{object}	MessageResponseType		"Failed to return tournament contest"
//	@Router			/tournament/{tournamentId}/contest [get]
func GetTournamentContest(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")
	if tournamentId == "" {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not a valid tournament id", tournamentId))
	}

	contestType := c.Query("type")
	if !models.CheckIfAllowedTournamentType(contestType) {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s is not allowed tournament format", contestType),
		)
	}
	tiktoks, err := database.GetTournamentTiktoksById(tournamentId)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("Could not get tiktoks for tournament with id %s", tournamentId))
	}
	shuffleTiktok(tiktoks)
	if contestType == models.SingleElimination {
		return c.Status(fiber.StatusOK).JSON(SingleElimination(tiktoks))
	}
	if contestType == models.KingOfTheHill {
		return c.Status(fiber.StatusOK).JSON(KingOfTheHill(tiktoks))
	}
	return MessageResponse(c, fiber.StatusBadRequest, "Unknown error")
}

func shuffleTiktok(t []models.Tiktok) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(t), func(i, j int) { t[i], t[j] = t[j], t[i] })
}

// SingleElimination https://en.wikipedia.org/wiki/Single-elimination_tournament
func SingleElimination(t []models.Tiktok) models.Bracket {
	countTiktok := len(t)
	countRound := int(math.Ceil(math.Log2(float64(countTiktok))))
	countSecondRoundParticipators := 1 << (countRound - 1) // Equivalent to int(math.Pow(2, float64(countRound)) / 2)
	countFirstRoundMatches := countTiktok - int(math.Pow(2, float64(countRound)-1))
	countFirstRoundTiktoks := countFirstRoundMatches * 2

	rounds := make([]models.Round, 0, countRound)

	firstRoundMatches := make([]models.Match, 0, countFirstRoundMatches)
	secondRoundMatches := make([]models.Match, 0, countSecondRoundParticipators/2)

	secondRoundParticipators := make([]models.Option, 0, countSecondRoundParticipators) // This slice should store MatchOption or TiktokOption

	// Filling first round with firstRoundMatches and appending MatchOptions to second round participators
	for j := 0; j < countFirstRoundTiktoks; j += 2 {
		matchID := uuid.NewString()
		firstRoundMatches = append(firstRoundMatches, models.Match{
			MatchID: matchID,
			FirstOption: models.TiktokOption{
				TiktokURL: t[j].URL,
			},
			SecondOption: models.TiktokOption{
				TiktokURL: t[j+1].URL,
			},
		})
		secondRoundParticipators = append(secondRoundParticipators,
			models.MatchOption{MatchID: matchID})
	}
	// Appending first round firstRoundMatches to rounds
	rounds = append(rounds, models.Round{
		Round:   1,
		Matches: firstRoundMatches,
	})
	// Appending TiktokOptions to second round participators
	for _, tiktok := range t[countFirstRoundTiktoks:] {
		secondRoundParticipators = append(secondRoundParticipators,
			models.TiktokOption{TiktokURL: tiktok.URL})
	}
	// Generating second round firstRoundMatches
	for i := 0; i < int(countSecondRoundParticipators); i += 2 {
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
				MatchID: uuid.NewString(),
				FirstOption: models.MatchOption{
					MatchID: previousRoundMatches[matchID].MatchID,
				},
				SecondOption: models.MatchOption{
					MatchID: previousRoundMatches[matchID+1].MatchID,
				},
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
	return models.Bracket{
		CountMatches: countTiktok - 1,
		Rounds:       rounds,
	}

}

// KingOfTheHill
// First match decided randomly between two participators.
// Loser of match leaves the game, winner will go to next match, next opponent decided randomly from standings.
// Procedure continues until last standing.
func KingOfTheHill(t []models.Tiktok) models.Bracket {
	countTiktok := len(t)
	rounds := make([]models.Round, 0, countTiktok-1)
	match := models.Match{
		MatchID:      uuid.NewString(),
		FirstOption:  models.TiktokOption{TiktokURL: t[0].URL},
		SecondOption: models.TiktokOption{TiktokURL: t[1].URL},
	}
	rounds = append(rounds, models.Round{
		Round:   1,
		Matches: []models.Match{match},
	})
	previousMatch := match
	for i := 2; i < countTiktok-1; i++ {
		match = models.Match{
			MatchID:      uuid.NewString(),
			FirstOption:  models.MatchOption{MatchID: previousMatch.MatchID},
			SecondOption: models.TiktokOption{TiktokURL: t[i].URL},
		}
		rounds = append(rounds, models.Round{
			Round:   i,
			Matches: []models.Match{match},
		})
		previousMatch = match
	}
	return models.Bracket{
		CountMatches: countTiktok - 1,
		Rounds:       rounds,
	}
}

// GetAllTournaments
//
//	@Summary		All tournaments
//	@Description	Get all tournaments
//	@Tags			tournament
//	@Accept			json
//	@Produce		json
//	@Success		200			{array}		models.Tournament	"Contest bracket"
//	@Failure		400			{object}	MessageResponseType	"Failed to return tournament contest"
//	@Router			/tournament																[get]
func GetAllTournaments(c *fiber.Ctx) error {
	tournaments, err := database.GetAllTournaments()
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, "Failed to get tournaments")
	}
	return c.Status(fiber.StatusOK).JSON(tournaments)
}
