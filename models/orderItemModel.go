package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Might need to fix the types later on.
//
//	The tutorial creator provided  deficient explanation as to why some are pointers and some are not.
type OrderItem struct {
	Id        primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Quantity  *int                `json:"quantity" bson:"quantity" validate:"required,eq=S|eq=M|eq=L"`
	UnitPrice *float64            `json:"unit_price" bson:"unit_price" validate:"required"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
	FoodId    *primitive.ObjectID `json:"food_id" bson:"food_id" validate:"required"`
	OrderId   primitive.ObjectID  `json:"order_id" bson:"order_id" validate:"required"`
}
