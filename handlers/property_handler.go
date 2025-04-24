package handlers

import (
	"dwello-api/db"
	"dwello-api/models"
	"dwello-api/utils"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetHomescreenProperties godoc
// @Summary Get properties for the homescreen
// @Description Get properties based on user's preferred location
// @Tags Properties
// @Accept json
// @Produce json
// @Success 200 {array} models.PropertySwagger
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/properties/homescreen [get]
func GetHomescreenProperties(c *fiber.Ctx) error {
	var user models.User
	userEmail := c.Query("email")
	fmt.Println("User email from query:", userEmail)
	if userEmail == "" {
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
		}
		userEmail = user.Email
	}

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	// Fetch user document
	err := db.UserCollection().FindOne(ctx, bson.M{"email": userEmail}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch user details"})
	}

	// Check for preferred locations
	if len(user.PreferredLocations) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Preferred locations not set"})
	}

	// Use $in to filter properties in any of the preferred locations
	filter := bson.M{"location": bson.M{"$in": user.PreferredLocations}}

	cursor, err := db.PropertyCollection().Find(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch properties"})
	}

	var properties []models.Property
	if err := cursor.All(ctx, &properties); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse properties"})
	}

	return c.JSON(properties)
}

// CreateProperty godoc
// @Summary Create a new property
// @Description Create a property owned by the authenticated user
// @Tags Properties
// @Accept json
// @Produce json
// @Param property body models.PropertySwagger true "Property data"
// @Success 201 {object} models.PropertySwagger
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/properties [post]
func CreateProperty(c *fiber.Ctx) error {
	// Parse only property-specific fields, not owner info
	var input struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Price       float64  `json:"price"`
		Location    string   `json:"location"`
		Thumbnail   string   `json:"thumbnail,omitempty"`
		Pictures    []string `json:"pictures,omitempty"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	// Find the user by email
	var user models.User
	if err := db.UserCollection().FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Create a property with the fetched user info
	property := models.Property{
		ID:          primitive.NewObjectID(),
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		Location:    input.Location,
		OwnerEmail:  user.Email,
		OwnerName:   user.Name,
		OwnerPic:    user.ProfilePic,
		IsRented:    false,
		Thumbnail:   input.Thumbnail,
		Pictures:    input.Pictures,
		CreatedAt:   primitive.NewDateTimeFromTime(utils.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(utils.Now()),
	}

	// Insert property into DB
	if _, err := db.PropertyCollection().InsertOne(ctx, property); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create property"})
	}

	// Add property ID to user's posted_properties
	_, err := db.UserCollection().UpdateOne(
		ctx,
		bson.M{"email": email},
		bson.M{"$push": bson.M{"posted_properties": property.ID}},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Property created, but failed to update user's posted properties"})
	}

	return c.Status(fiber.StatusCreated).JSON(property)
}

// UpdateProperty godoc
// @Summary Update an existing property
// @Description Update a property owned by the authenticated user
// @Tags Properties
// @Accept json
// @Produce json
// @Param id path string true "Property ID"
// @Param property body models.PropertySwagger true "Updated property data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/properties/{id} [put]
func UpdateProperty(c *fiber.Ctx) error {
	propertyID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid property ID"})
	}

	var property models.Property
	if err := c.BodyParser(&property); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Get email from body or query
	userEmail := property.OwnerEmail
	if userEmail == "" {
		userEmail = c.Query("email")
	}

	if userEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	// Check if property exists and belongs to the user
	var existingProperty models.Property
	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	err = db.PropertyCollection().FindOne(ctx, bson.M{"_id": propertyID}).Decode(&existingProperty)
	if err != nil || existingProperty.OwnerEmail != userEmail {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You cannot update a property that doesn't belong to you"})
	}

	property.UpdatedAt = primitive.NewDateTimeFromTime(utils.Now())

	_, err = db.PropertyCollection().UpdateOne(ctx, bson.M{"_id": propertyID}, bson.M{"$set": property})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update property"})
	}

	return c.JSON(fiber.Map{"message": "Property updated"})
}

// DeleteProperty godoc
// @Summary Delete a property
// @Description Delete a property owned by the authenticated user
// @Tags Properties
// @Accept json
// @Produce json
// @Param id path string true "Property ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/properties/{id} [delete]
func DeleteProperty(c *fiber.Ctx) error {
	propertyID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid property ID"})
	}

	// Get email from body or query
	var property models.Property
	userEmail := c.Query("email")
	if userEmail == "" {
		if err := c.BodyParser(&property); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
		}
		userEmail = property.OwnerEmail
	}

	if userEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	// Check if the property exists and belongs to the user
	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	err = db.PropertyCollection().FindOne(ctx, bson.M{"_id": propertyID}).Decode(&property)
	if err != nil || property.OwnerEmail != userEmail {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You cannot delete a property that doesn't belong to you"})
	}

	_, err = db.PropertyCollection().DeleteOne(ctx, bson.M{"_id": propertyID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete property"})
	}

	return c.JSON(fiber.Map{"message": "Property deleted"})
}

// LikeProperty godoc
// @Summary Like a property
// @Description Add a property to the user's liked list
// @Tags Properties
// @Accept json
// @Produce json
// @Param id path string true "Property ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/properties/{id}/like [post]
func LikeProperty(c *fiber.Ctx) error {
	propertyID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid property ID"})
	}

	// Get email from body or query
	userEmail := c.Query("email")
	if userEmail == "" {
		var property models.Property
		if err := c.BodyParser(&property); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
		}
		userEmail = property.OwnerEmail
	}

	if userEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	_, err = db.PropertyCollection().UpdateOne(ctx, bson.M{"_id": propertyID}, bson.M{"$addToSet": bson.M{"liked_by": userEmail}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to like property"})
	}

	_, err = db.UserCollection().UpdateOne(ctx, bson.M{"email": userEmail}, bson.M{"$addToSet": bson.M{"liked_properties": propertyID}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user liked properties"})
	}

	return c.JSON(fiber.Map{"message": "Property liked"})
}

// UnlikeProperty godoc
// @Summary Unlike a property
// @Description Remove a property from the user's liked list
// @Tags Properties
// @Accept json
// @Produce json
// @Param id path string true "Property ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/properties/{id}/unlike [post]
func UnlikeProperty(c *fiber.Ctx) error {
	propertyID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid property ID"})
	}

	// Get email from body or query
	userEmail := c.Query("email")
	if userEmail == "" {
		var property models.Property
		if err := c.BodyParser(&property); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
		}
		userEmail = property.OwnerEmail
	}

	if userEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	_, err = db.PropertyCollection().UpdateOne(ctx, bson.M{"_id": propertyID}, bson.M{"$pull": bson.M{"liked_by": userEmail}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to unlike property"})
	}

	_, err = db.UserCollection().UpdateOne(ctx, bson.M{"email": userEmail}, bson.M{"$pull": bson.M{"liked_properties": propertyID}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user liked properties"})
	}

	return c.JSON(fiber.Map{"message": "Property unliked"})
}

// GetUserLikedProperties godoc
// @Summary Get properties liked by the user
// @Description Get properties liked by the user based on email
// @Tags Properties
// @Accept json
// @Produce json
// @Param email query string true "User email"
// @Success 200 {array} models.PropertySwagger
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/properties/liked-properties [get]
func GetLikedPropertiesByUser(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	// Step 1: Find the user by email
	var user models.User
	if err := db.UserCollection().FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Step 2: If the user has no liked properties
	if len(user.LikedProperties) == 0 {
		return c.JSON([]models.Property{}) // return empty array
	}

	// Step 3: Fetch the liked properties
	cursor, err := db.PropertyCollection().Find(ctx, bson.M{
		"_id": bson.M{"$in": user.LikedProperties},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch liked properties"})
	}

	var properties []models.Property
	if err := cursor.All(ctx, &properties); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode properties"})
	}

	return c.JSON(properties)
}

// SearchProperties godoc
// @Summary Search properties
// @Description Search for properties by location, price, etc.
// @Tags Properties
// @Accept json
// @Produce json
// @Param location query string false "Location"
// @Param min_price query string false "Minimum price"
// @Param max_price query string false "Maximum price"
// @Param limit query int false "Limit"
// @Param skip query int false "Skip"
// @Success 200 {array} models.PropertySwagger
// @Failure 500 {object} map[string]string
// @Router /api/properties/search [get]
func SearchProperties(c *fiber.Ctx) error {
	location := c.Query("location")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")

	filter := bson.M{}
	if location != "" {
		filter["location"] = location
	}

	if minPrice != "" || maxPrice != "" {
		priceFilter := bson.M{}
		if minVal, err := strconv.Atoi(minPrice); err == nil {
			priceFilter["$gte"] = minVal
		}
		if maxVal, err := strconv.Atoi(maxPrice); err == nil {
			priceFilter["$lte"] = maxVal
		}
		if len(priceFilter) > 0 {
			filter["price"] = priceFilter
		}
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	skip, _ := strconv.Atoi(c.Query("skip", "0"))

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	cursor, err := db.PropertyCollection().Find(ctx, filter, options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch properties"})
	}

	var properties []models.Property
	if err := cursor.All(ctx, &properties); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse properties"})
	}
	return c.JSON(properties)
}

// RequestToRentProperty godoc
// @Summary Request to rent a property
// @Description Send a rental request for a property
// @Tags Properties
// @Accept json
// @Produce json
// @Param id path string true "Property ID"
// @Param email query string true "User email"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/properties/{id}/request [post]
func RequestToRentProperty(c *fiber.Ctx) error {
	propertyIDParam := c.Params("id")
	userIDParam := c.Query("user_id")

	if propertyIDParam == "" || userIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Property ID and user ID are required"})
	}

	propertyID, err := primitive.ObjectIDFromHex(propertyIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid property ID"})
	}

	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	ctx, cancel := utils.DatabaseContext()
	defer cancel()

	// Add user ID to property's rental_requests
	_, err = db.PropertyCollection().UpdateOne(ctx, bson.M{"_id": propertyID}, bson.M{
		"$addToSet": bson.M{"rental_requests": userID},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add rental request"})
	}

	// Add property ID to user's rental_requests
	_, err = db.UserCollection().UpdateOne(ctx, bson.M{"_id": userID}, bson.M{
		"$addToSet": bson.M{"rental_requests": propertyID},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user's request list"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Rental request sent"})
}
