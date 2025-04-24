package handlers

import (
	"dwello-api/db"
	"dwello-api/models"
	"dwello-api/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterUser registers a new user or logs in if the user already exists
// @Summary Register or Login User
// @Description Register a new user or return existing user if already registered
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.UserSwagger true "User JSON"
// @Success 200 {object} models.UserSwagger
// @Success 201 {object} models.UserSwagger
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/register [post]
func RegisterUser(c *fiber.Ctx) error {
	fmt.Println("RegisterUser called")
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		// return what is missing in the body
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body", "details": err.Error()})
	}

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	collection := db.UserCollection()

	// Check if user already exists
	var existing models.User
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existing)
	if err == nil {
		return c.Status(fiber.StatusOK).JSON(existing) // User already exists, return it
	}

	// New user
	user.ID = primitive.NewObjectID()
	user.PostedProperties = []primitive.ObjectID{}
	user.LikedProperties = []primitive.ObjectID{}
	user.CreatedAt = primitive.NewDateTimeFromTime(utils.Now())
	user.UpdatedAt = user.CreatedAt

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUserByEmail fetches a user by their email
// @Summary Get User by Email
// @Description Get a user document based on email address
// @Tags Users
// @Produce json
// @Param email path string true "User Email"
// @Success 200 {object} models.UserSwagger
// @Failure 404 {object} map[string]string
// @Router /api/users/{email} [get]
func GetUserByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	var user models.User
	err := db.UserCollection().FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

// UpdateUserLocation updates the current location of a user
// @Summary Update User Location
// @Description Update the location field of a user
// @Tags Users
// @Accept json
// @Produce json
// @Param email path string true "User Email"
// @Param location body map[string]string true "Location JSON"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{email}/location [put]
func UpdateUserLocation(c *fiber.Ctx) error {
	email := c.Params("email")
	var payload struct {
		Location string `json:"location"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	_, err := db.UserCollection().UpdateOne(ctx,
		bson.M{"email": email},
		bson.M{
			"$set": bson.M{
				"location":   payload.Location,
				"updated_at": primitive.NewDateTimeFromTime(utils.Now()),
			},
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}
	return c.JSON(fiber.Map{"message": "Location updated"})
}

// UpdatePreferredLocations updates the user's preferred locations
// @Summary Update Preferred Locations
// @Description Update the preferred_locations field of a user (can be multiple)
// @Tags Users
// @Accept json
// @Produce json
// @Param email path string true "User Email"
// @Param preferred_locations body map[string][]string true "Preferred Locations JSON"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{email}/preferred-locations [put]
func UpdatePreferredLocations(c *fiber.Ctx) error {
	email := c.Params("email")
	var payload struct {
		PreferredLocations []string `json:"preferred_locations"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	_, err := db.UserCollection().UpdateOne(ctx,
		bson.M{"email": email},
		bson.M{
			"$set": bson.M{
				"preferred_locations": payload.PreferredLocations,
				"updated_at":          primitive.NewDateTimeFromTime(utils.Now()),
			},
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}
	return c.JSON(fiber.Map{"message": "Preferred locations updated"})
}

// GetLikedProperties retrieves the properties liked by a user
// @Summary Get Liked Properties
// @Description Get a list of properties the user has liked
// @Tags Users
// @Produce json
// @Param email path string true "User Email"
// @Success 200 {array} models.PropertySwagger
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{email}/liked-properties [get]
func GetLikedProperties(c *fiber.Ctx) error {
	email := c.Params("email")

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	// Get user
	var user models.User
	err := db.UserCollection().FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Get properties by IDs
	cursor, err := db.PropertyCollection().Find(ctx, bson.M{
		"_id": bson.M{"$in": user.LikedProperties},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch properties"})
	}

	var properties []models.Property
	if err := cursor.All(ctx, &properties); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse properties"})
	}
	return c.JSON(properties)
}

// GetPostedProperties retrieves the properties posted by a user
// @Summary Get Posted Properties
// @Description Get a list of properties the user has posted
// @Tags Users
// @Produce json
// @Param email path string true "User Email"
// @Success 200 {array} models.PropertySwagger
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{email}/posted-properties [get]
func GetPostedProperties(c *fiber.Ctx) error {
	email := c.Params("email")

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	// Get user
	var user models.User
	err := db.UserCollection().FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Get properties by IDs
	cursor, err := db.PropertyCollection().Find(ctx, bson.M{
		"_id": bson.M{"$in": user.PostedProperties},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch properties"})
	}

	var properties []models.Property
	if err := cursor.All(ctx, &properties); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse properties"})
	}
	return c.JSON(properties)
}

// HandleRentalRequest handles rental requests from users
// @Summary Handle Rental Request
// @Description Accept or reject a rental request for a property
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "Property ID"
// @Param renter query string true "Renter Email"
// @Param action query string true "Action (accept/reject)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/rental-request/{id} [post]
func HandleRentalRequest(c *fiber.Ctx) error {
	propertyIDParam := c.Params("id")
	renterIDParam := c.Query("renter_id")
	action := c.Query("action") // "accept" or "reject"

	if propertyIDParam == "" || renterIDParam == "" || action == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Property ID, renter ID, and action are required"})
	}

	propertyID, err := primitive.ObjectIDFromHex(propertyIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid property ID"})
	}

	renterID, err := primitive.ObjectIDFromHex(renterIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid renter ID"})
	}

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	// Always remove the request from both sides
	_, err = db.PropertyCollection().UpdateOne(ctx, bson.M{"_id": propertyID}, bson.M{
		"$pull": bson.M{"rental_requests": renterID},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update property"})
	}

	_, err = db.UserCollection().UpdateOne(ctx, bson.M{"_id": renterID}, bson.M{
		"$pull": bson.M{"rental_requests": propertyID},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	if action == "accept" {
		// Mark the property as rented
		_, err = db.PropertyCollection().UpdateOne(ctx, bson.M{"_id": propertyID}, bson.M{
			"$set": bson.M{
				"is_rented":    true,
				"rented_by_id": renterID,
			},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to mark property as rented"})
		}

		// Add to user's rented properties
		_, err = db.UserCollection().UpdateOne(ctx, bson.M{"_id": renterID}, bson.M{
			"$addToSet": bson.M{"rented_properties": propertyID},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user rented properties"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Rental request accepted"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Rental request rejected"})
}

// GetRentalRequestsForUserProperties retrieves rental requests for properties owned by a user
// @Summary Get Rental Requests for User Properties
// @Description Get rental requests for properties owned by a user
// @Tags Users
// @Produce json
// @Param email query string true "User Email"
// @Success 200 {array} models.PropertySwagger
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/rental-requests [get]
func GetRentalRequestsForUserProperties(c *fiber.Ctx) error {
	userEmail := c.Query("email")
	if userEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"owner_email":     userEmail,
			"rental_requests": bson.M{"$ne": bson.A{}},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "rental_requests",
			"foreignField": "_id",
			"as":           "requesting_users",
		}}},
	}

	cursor, err := db.PropertyCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch rental requests"})
	}

	var properties []bson.M
	if err := cursor.All(ctx, &properties); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode rental requests"})
	}

	return c.JSON(properties)
}

// GetRentedPropertiesByUser retrieves properties rented by a user
// @Summary Get Rented Properties by User
// @Description Get properties rented by a user
// @Tags Users
// @Produce json
// @Param email query string true "User Email"
// @Success 200 {array} models.PropertySwagger
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/rented-properties [get]
func GetRentedPropertiesByUser(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID required"})
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid User ID"})
	}

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	cursor, err := db.PropertyCollection().Find(ctx, bson.M{"rented_by_id": objectID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch rented properties"})
	}

	var properties []models.Property
	if err := cursor.All(ctx, &properties); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode properties"})
	}

	return c.JSON(properties)
}
