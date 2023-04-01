package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Might need to fix the types later on.
//
//	The tutorial creator provided  deficient explanation as to why some are pointers and some are not.
type User struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName    *string            `json:"first_name" bson:"first_name" validate:"required,min=2,max=100"`
	LastName     *string            `json:"last_name" bson:"last_name" validate:"required,min=2,max=100"`
	Password     *string            `json:"password" bson:"password" validate:"required,min=8,max=100"`
	Email        *string            `json:"email" bson:"email" validate:"required,email"`
	Avatar       *string            `json:"avatar" bson:"avatar"`
	Phone        *string            `json:"phone" bson:"phone" validate:"required"`
	Token        *string            `json:"token" bson:"token"`
	RefreshToken *string            `json:"refresh_token" bson:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
