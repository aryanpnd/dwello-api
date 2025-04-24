package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                 primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Email              string               `bson:"email" json:"email"`
	Name               string               `bson:"name" json:"name"`
	ProfilePic         string               `bson:"profile_pic,omitempty" json:"profile_pic,omitempty"`
	Location           string               `bson:"location,omitempty" json:"location,omitempty"`
	PreferredLocations []string             `bson:"preferred_locations,omitempty" json:"preferred_locations,omitempty"`
	PostedProperties   []primitive.ObjectID `bson:"posted_properties,omitempty" json:"posted_properties,omitempty"`
	LikedProperties    []primitive.ObjectID `bson:"liked_properties,omitempty" json:"liked_properties,omitempty"`
	RentedProperties   []primitive.ObjectID `bson:"rented_properties,omitempty" json:"rented_properties,omitempty"`
	RentalRequests     []primitive.ObjectID `bson:"rental_requests,omitempty" json:"rental_requests,omitempty"` // properties the user has requested

	CreatedAt primitive.DateTime `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt primitive.DateTime `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// UserSwagger is a Swagger-friendly version of User
type UserSwagger struct {
	Email              string   `json:"email" example:"user@example.com"`
	Name               string   `json:"name" example:"Alice Smith"`
	ProfilePic         string   `json:"profile_pic,omitempty"`
	Location           string   `json:"location,omitempty"`
	PreferredLocations []string `json:"preferred_locations,omitempty" example:"[\"Los Angeles\", \"New York\"]"`
	PostedProperties   []string `json:"posted_properties,omitempty"`
	LikedProperties    []string `json:"liked_properties,omitempty"`
	RentedProperties   []string `json:"rented_properties,omitempty"`
	RentalRequests     []string `json:"rental_requests,omitempty"` // properties the user has requested
}
