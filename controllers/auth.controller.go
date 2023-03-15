package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"tiktok-arena/configuration"
	"tiktok-arena/database"
	"tiktok-arena/models"
	"time"
)

// RegisterUser
//	@Summary		Register user
//	@Description	Register new user with given credentials
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload			body		models.RegisterInput	true	"Data to register user"
//	@Success		200				{object}	models.UserAuthDetails	"Register success"
//	@Failure		400				{object}	MessageResponseType		"Failed to register user"
//	@Router			/auth/register	[post]
func RegisterUser(c *fiber.Ctx) error {
	var payload *models.RegisterInput

	err := c.BodyParser(&payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = models.ValidateStruct(payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if database.CheckIfUserExists(payload.Name) {
		return MessageResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("User %s already exists", payload.Name))
	}

	newUser := models.User{
		Name:     payload.Name,
		Password: string(hashedPassword),
	}

	err = database.CreateNewUser(&newUser)

	if err != nil {
		return MessageResponse(c, fiber.StatusBadGateway, err.Error())
	}

	token, err := UserJwtToken(&newUser)

	if err != nil {
		return MessageResponse(c, fiber.StatusBadGateway,
			fmt.Sprintf("Generating JWT Token failed: %v", err))
	}

	return c.Status(fiber.StatusCreated).JSON(
		models.UserAuthDetails{
			ID:       newUser.ID.String(),
			Username: newUser.Name,
			Token:    token,
		},
	)
}

// LoginUser
//	@Summary		Login user
//	@Description	Login user with given credentials
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload			body		models.LoginInput		true	"Data to login user"
//	@Success		200				{object}	models.UserAuthDetails	"Login success"
//	@Failure		400				{object}	MessageResponseType		"Error logging in"
//	@Router			/auth/login    	[post]
func LoginUser(c *fiber.Ctx) error {
	var payload *models.LoginInput

	err := c.BodyParser(&payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = models.ValidateStruct(payload)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, err.Error())
	}

	user, err := database.GetUserByName(payload.Name)
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, "Invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return MessageResponse(c, fiber.StatusBadRequest, "Invalid credentials")
	}

	token, err := UserJwtToken(&user)

	if err != nil {
		return MessageResponse(c, fiber.StatusBadGateway,
			fmt.Sprintf("Generating JWT Token failed: %v", err))
	}

	return c.Status(fiber.StatusOK).JSON(
		models.UserAuthDetails{
			ID:       user.ID.String(),
			Username: user.Name,
			Token:    token,
		},
	)
}

func UserJwtToken(user *models.User) (string, error) {
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
		return "", err
	}

	return tokenString, nil
}

// WhoAmI
//	@Summary		Authenticated user details
//	@Description	Get current user id and name
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	models.UserAuthDetails	"User details"
//	@Failure		400	{object}	MessageResponseType		"Error getting user data"
//	@Router			/auth/whoami [get]
func WhoAmI(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	username := claims["name"].(string)
	id := claims["sub"].(string)

	return c.Status(fiber.StatusOK).JSON(models.UserAuthDetails{
		ID:       id,
		Username: username,
		Token:    token.Raw,
	})
}
