package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"tiktok-arena/configuration"
	"tiktok-arena/database"
	"tiktok-arena/router"
)

func init() {
	err := configuration.LoadConfig(".env")
	if err != nil {
		log.Fatalln("Failed to load environment variables!", err.Error())
	}
	database.ConnectDB(&configuration.EnvConfig)
}

//	@title			TikTok arena API
//	@version		1.0
//	@description	API for TikTok arena application
//	@host			localhost:8000
//	@BasePath		/api
func main() {
	app := fiber.New()

	//	Logger middleware for logging HTTP request/response details
	app.Use(logger.New())

	//	Cors middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":8000"))
}
