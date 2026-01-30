package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	FirstName    string `bson:"first_name" json:"first_name"`
	LastName     string `bson:"last_name" json:"last_name"`
	Email        string `bson:"email" json:"email"`
	PasswordHash string `bson:"password_hash" json:"-"`
	Role         string `bson:"role" json:"role"`

	Balance float64 `bson:"balance" json:"balance"`

	Watchlist []PriceAlert    `bson:"watchlist" json:"watchlist"`
	Portfolio []BoughtStock  `bson:"portfolio" json:"portfolio"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
