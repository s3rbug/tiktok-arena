package router

import (
	"github.com/gofiber/fiber/v2"
	"tiktok-arena/controllers"
	"tiktok-arena/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Route("/auth", func(router fiber.Router) {
		router.Post("/register", controllers.RegisterUser)
		router.Post("/login", controllers.LoginUser)
		router.Get("/whoami", middleware.Protected(), controllers.WhoAmI)
	})
}
