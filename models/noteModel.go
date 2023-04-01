package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Might need to fix the types later on.
//
//	The tutorial creator provided  deficient explanation as to why some are pointers and some are not.
type Note struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Text      string             `json:"text" bson:"text" validate:"required,min=2,max=100"`
	Title     string             `json:"title" bson:"title" validate:"required,min=2,max=100"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
