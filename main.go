package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"tiktok-arena/configuration"
	"tiktok-arena/database"
	"tiktok-arena/router"
)

func init() {
	err := configuration.LoadConfig(".env")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	database.ConnectDB(&configuration.EnvConfig)
}

func main() {
	app := fiber.New()
	//	Logger middleware for logging HTTP request/response details
	app.Use(logger.New())

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":8000"))
}
