package services

import (
	"context"
	"strings"
	"time"

	db "github.com/GeorgiStoyanov05/GoMarket2/database"
	"github.com/GeorgiStoyanov05/GoMarket2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListUserPositions(userID primitive.ObjectID) ([]models.Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	coll := db.Client.Database("gomarket").Collection("positions")

	cur, err := coll.Find(ctx, bson.M{"user_id": userID, "qty": bson.M{"$gt": 0}},
		options.Find().SetSort(bson.D{{Key: "symbol", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]models.Position, 0)
	for cur.Next(ctx) {
		var p models.Position
		if err := cur.Decode(&p); err != nil {
			continue
		}
		p.Symbol = strings.ToUpper(strings.TrimSpace(p.Symbol))
		out = append(out, p)
	}
	return out, nil
}
