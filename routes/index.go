package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Mount route groups
	RegisterUserRoutes(app)
	RegisterPropertyRoutes(app)
}
