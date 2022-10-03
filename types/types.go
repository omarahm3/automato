package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Hash      string             `json:"hash" bson:"hash"`
	Title     string             `json:"title" bson:"title"`
	Video     string             `json:"video" bson:"video"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
