package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Might need to fix the types later on.
//
//	The tutorial creator provided  deficient explanation as to why some are pointers and some are not.
type Food struct {
	Id        primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Name      *string             `json:"name" bson:"name" validate:"required,min=2,max=100"`
	Price     *float64            `json:"price" bson:"price" validate:"required`
	Image     *string             `json:"image" bson:"image" validate:"required`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
	MenuId    *primitive.ObjectID `json:"menu_id" bson:"menu_id" validate:"required`
}
