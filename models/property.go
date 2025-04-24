package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Property struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Location    string             `bson:"location" json:"location"`

	OwnerEmail string `bson:"owner_email" json:"owner_email"`
	OwnerName  string `bson:"owner_name" json:"owner_name"`
	OwnerPic   string `bson:"owner_pic" json:"owner_pic"`

	IsRented       bool                 `bson:"is_rented" json:"is_rented"`
	RentedByEmail  string               `bson:"rented_by_email,omitempty" json:"rented_by_email,omitempty"`
	RentalRequests []primitive.ObjectID `bson:"rental_requests,omitempty" json:"rental_requests,omitempty"`

	Thumbnail string   `bson:"thumbnail,omitempty" json:"thumbnail,omitempty"`
	Pictures  []string `bson:"pictures,omitempty" json:"pictures,omitempty"`

	LikedBy   []string           `bson:"liked_by,omitempty" json:"liked_by,omitempty"`
	CreatedAt primitive.DateTime `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt primitive.DateTime `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// PropertySwagger is a Swagger-friendly version of Property
type PropertySwagger struct {
	Title       string  `json:"title" example:"Modern 2BHK Apartment"`
	Description string  `json:"description" example:"Spacious apartment near downtown."`
	Price       float64 `json:"price" example:"2500"`
	Location    string  `json:"location" example:"New York"`

	OwnerEmail string `json:"owner_email" example:"owner@example.com"`
	OwnerName  string `json:"owner_name" example:"John Doe"`
	OwnerPic   string `json:"owner_pic" example:"https://example.com/pic.jpg"`

	RentedBy       []string `json:"rented_by,omitempty"`
	RentalRequests []string `json:"rental_requests,omitempty"`
	IsRented       bool     `json:"is_rented" example:"false"`

	Thumbnail string   `json:"thumbnail,omitempty"`
	Pictures  []string `json:"pictures,omitempty"`

	LikedBy []string `json:"liked_by,omitempty"`
}
