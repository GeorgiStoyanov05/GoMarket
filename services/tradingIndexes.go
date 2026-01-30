package services

import (
	"context"
	"time"

	db "github.com/GeorgiStoyanov05/GoMarket/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureTradingIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d := db.Client.Database("gomarket")

	positions := d.Collection("positions")
	orders := d.Collection("orders")

	// One position per (user, symbol)
	_, _ = positions.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "symbol", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// Helpful for listing user orders
	_, _ = orders.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: -1}},
	})
}
