package db

import (
	"dwello-api/config"

	"go.mongodb.org/mongo-driver/mongo"
)

func UserCollection() *mongo.Collection {
	return config.DB.Collection("users")
}

func PropertyCollection() *mongo.Collection {
	return config.DB.Collection("properties")
}
