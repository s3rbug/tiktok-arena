package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"tiktok-arena/database"
	"tiktok-arena/models"
)

// TournamentsOfUser
//
//	@Summary		Create new tournament
//	@Description	Create new tournament for current user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		string						false	"page number"
//	@Param			count		query		string						false	"page size"
//	@Param			sort_name	query		string						false	"sort page by name"
//	@Param			sort_size	query		string						false	"sort page by size"
//	@Param			search		query		string						false	"search"
//	@Success		200			{object}	models.TournamentsResponse	"Tournaments of user"
//	@Failure		400			{object}	MessageResponseType			"Couldn't get tournaments for specific user"
//	@Router			/user/tournaments [get]
func TournamentsOfUser(c *fiber.Ctx) error {
	userId, err := GetUserIdAndCheckJWT(c)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}
	p := new(models.PaginationQueries)
	if err := c.QueryParser(p); err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, "Failed to parse queries")
	}
	models.ValidatePaginationQueries(p)
	tournamentResponse, err := database.GetAllTournamentsForUserById(userId, *p)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, "Failed to get tournaments")
	}

	if tournamentResponse.TournamentCount == 0 {
		return MessageResponse(c, fiber.StatusMovedPermanently, "There is no page")
	}

	return c.Status(fiber.StatusOK).JSON(tournamentResponse)
}

func GetUserIdAndCheckJWT(c *fiber.Ctx) (uuid.UUID, error) {
	user := c.Locals("user")

	if user == nil {
		return uuid.UUID{}, MessageResponse(c, fiber.StatusBadRequest, "Empty jwt.token")
	}
	userJWT := user.(*jwt.Token)

	claims := userJWT.Claims.(jwt.MapClaims)

	userId, err := uuid.Parse(claims["sub"].(string))

	return userId, err
}
