package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Position struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`

	Symbol  string `bson:"symbol" json:"symbol"`
	Qty     int64  `bson:"qty" json:"qty"`
	AvgCost float64 `bson:"avg_cost" json:"avg_cost"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
