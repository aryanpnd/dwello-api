package routes

import (
	"dwello-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterPropertyRoutes(app *fiber.App) {
	// Grouping the property-related routes
	property := app.Group("/api/properties")

	// Create a new property
	property.Post("/", handlers.CreateProperty)

	// Update an existing property
	property.Put("/:id", handlers.UpdateProperty)

	// Delete a property
	property.Delete("/:id", handlers.DeleteProperty)

	// Like a property
	property.Post("/:id/like", handlers.LikeProperty)

	// Unlike a property
	property.Post("/:id/unlike", handlers.UnlikeProperty)

	// Get user liked properties
	property.Get("/liked-properties", handlers.GetLikedPropertiesByUser)

	// Search for properties
	property.Get("/search", handlers.SearchProperties)

	// Get properties for the homescreen based on preferred location
	property.Get("/homescreen", handlers.GetHomescreenProperties)

	// Rental features
	property.Post("/:id/rent", handlers.RequestToRentProperty)
}
