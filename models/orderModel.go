package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Might need to fix the types later on.
//
//	The tutorial creator provided  deficient explanation as to why some are pointers and some are not.
type Order struct {
	Id        primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	OrderDate time.Time           `json:"order_date" bson:"order_date" validate:"required"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
	TableId   *primitive.ObjectID `json:"table_id" bson:"table_id" validate:"required"`
}
