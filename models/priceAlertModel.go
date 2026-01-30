package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PriceAlert struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Symbol    string             `bson:"symbol" json:"symbol"`
	Condition string             `bson:"condition" json:"condition"`

	TargetPrice float64 `bson:"target_price" json:"target_price"`

	Active         bool      `bson:"active" json:"active"`
	Triggered      bool      `bson:"triggered" json:"triggered"`
	TriggeredAt    time.Time `bson:"triggered_at" json:"triggered_at"`
	TriggeredPrice float64   `bson:"triggered_price" json:"triggered_price"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
