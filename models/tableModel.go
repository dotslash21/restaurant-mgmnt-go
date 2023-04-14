package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Might need to fix the types later on.
//
//	The tutorial creator provided  deficient explanation as to why some are pointers and some are not.
type Table struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	NumberOfGuests *int               `json:"number_of_guests" bson:"number_of_guests" validate:"required"`
	Number         *int               `json:"number" bson:"number" validate:"required"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}
