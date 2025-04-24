package routes

import (
	"dwello-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App) {
	// Grouping the user-related routes
	user := app.Group("/api/users")

	// Register a new user or login
	user.Post("/register", handlers.RegisterUser)

	// Get user details by email
	user.Get("/:email", handlers.GetUserByEmail)

	// Update user location
	user.Put("/:email/location", handlers.UpdateUserLocation)

	// Update user preferred locations
	user.Put("/:email/preferred-locations", handlers.UpdatePreferredLocations)

	// Get properties liked by the user
	user.Get("/:email/liked-properties", handlers.GetLikedProperties)

	// Get properties posted by the user
	user.Get("/:email/posted-properties", handlers.GetPostedProperties)

	// Get properties rented by this user
	user.Get("/:email/rented-properties", handlers.GetRentedPropertiesByUser)

	// Get rental requests received for this user's properties
	user.Get("/:email/rental-requests", handlers.GetRentalRequestsForUserProperties)

	// Handle rental request for a property
	user.Post("/rental-requests/:id/handle", handlers.HandleRentalRequest)

}
