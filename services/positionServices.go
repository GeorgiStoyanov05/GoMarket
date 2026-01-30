package services

import (
	"context"
	"strings"
	"time"

	db "github.com/GeorgiStoyanov05/GoMarket2/database"
	"github.com/GeorgiStoyanov05/GoMarket2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserPosition(userID primitive.ObjectID, symbol string) (*models.Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	sym := strings.ToUpper(strings.TrimSpace(symbol))
	coll := db.Client.Database("gomarket").Collection("positions")

	var p models.Position
	err := coll.FindOne(ctx, bson.M{"user_id": userID, "symbol": sym}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
