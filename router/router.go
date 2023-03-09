package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"tiktok-arena/controllers"
	_ "tiktok-arena/docs"
	"tiktok-arena/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	//	Use 'swag init' to generate new /docs files, details: https://github.com/gofiber/swagger#usage
	api.Get("/docs/*", swagger.HandlerDefault)

	api.Route("/auth", func(router fiber.Router) {
		router.Post("/register", controllers.RegisterUser)
		router.Post("/login", controllers.LoginUser)
		router.Get("/whoami", middleware.Protected(), controllers.WhoAmI)
	})

	api.Route("/tournament", func(router fiber.Router) {
		router.Post("", middleware.Protected(), controllers.CreateTournament)
		router.Get("/:tournamentId", controllers.GetTournamentDetails)
		router.Get("/:tournamentId/tiktoks", controllers.GetTournamentTiktoks)
	})
}
