package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Might need to fix the types later on.
//
//	The tutorial creator provided  deficient explanation as to why some are pointers and some are not.
type Invoice struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OrderId        primitive.ObjectID `json:"order_id" bson:"order_id" validate:"required"`
	PaymentMethod  *string            `json:"payment_method" bson:"payment_method" validate:"eq=CARD|eq=CASH|eq="`
	PaymentStatus  *string            `json:"payment_status" bson:"payment_status" validate:"required,eq=PAID|eq=UNPAID"`
	PaymentDueDate time.Time          `json:"payment_due_date" bson:"payment_due_date"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}
