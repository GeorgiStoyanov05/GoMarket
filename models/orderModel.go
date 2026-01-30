package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`

	Symbol string `bson:"symbol" json:"symbol"`
	Side   string `bson:"side" json:"side"` // "buy" | "sell"

	Qty   int64   `bson:"qty" json:"qty"`
	Price float64 `bson:"price" json:"price"` // fill price (market = quote at time)

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
