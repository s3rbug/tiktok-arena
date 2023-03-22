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
		router.Get("", controllers.GetAllTournaments)
		router.Get("/contest/:tournamentId", controllers.GetTournamentContest)
		router.Post("/create", middleware.Protected(), controllers.CreateTournament)
		router.Post("/edit/:tournamentId", middleware.Protected(), controllers.EditTournament)
		router.Delete("/delete/:tournamentId", middleware.Protected(), controllers.DeleteTournament)
		router.Delete("/delete", middleware.Protected(), controllers.DeleteTournaments)
		router.Get("/tiktoks/:tournamentId", controllers.GetTournamentTiktoks)
		router.Get("/:tournamentId", controllers.GetTournamentDetails)
		router.Post("/:tournamentId", controllers.TournamentWinner)
	})

	api.Route("/user", func(router fiber.Router) {
		router.Get("/tournaments", middleware.Protected(), controllers.TournamentsOfUser)
	})
}
