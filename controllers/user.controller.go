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
//	@Success		200	{object}	MessageResponseType	"Tournaments of user"
//	@Failure		400	{object}	MessageResponseType	"Couldn't get tournaments for specific user"
//	@Router			/user/tournaments [get]
func TournamentsOfUser(c *fiber.Ctx) error {
	userId, err := GetUserIdAndCheckJWT(c)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}
	var tournamentsOfUser []models.Tournament
	tournamentsOfUser, err = database.GetAllTournamentsForUserById(userId)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, "Failed to get tournaments")
	}
	return c.Status(fiber.StatusOK).JSON(tournamentsOfUser)
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
