package services

import (
	"context"
	"math"
	"strings"
	"time"

	db "github.com/GeorgiStoyanov05/GoMarket/database"
	"github.com/GeorgiStoyanov05/GoMarket/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const alertsCollection = "alerts"

func roundToCents(v float64) float64 {
	return math.Round(v*100) / 100
}

func CreatePriceAlert(userID primitive.ObjectID, symbol, condition string, targetPrice float64) (models.PriceAlert, map[string]string) {
	errs := map[string]string{}

	sym := strings.ToUpper(strings.TrimSpace(symbol))
	cond := strings.ToLower(strings.TrimSpace(condition))
	target := roundToCents(targetPrice)

	if sym == "" {
		errs["symbol"] = "Missing symbol."
	}
	if cond != "above" && cond != "below" {
		errs["condition"] = "Condition must be 'above' or 'below'."
	}
	if target <= 0 {
		errs["targetPrice"] = "Target price must be bigger than 0."
	}
	if len(errs) > 0 {
		return models.PriceAlert{}, errs
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	coll := db.Client.Database("gomarket").Collection(alertsCollection)

	// Basic idempotency: block exact duplicates (same user + symbol + condition + target, still active).
	var existing models.PriceAlert
	err := coll.FindOne(ctx, bson.M{
		"user_id":      userID,
		"symbol":       sym,
		"condition":    cond,
		"target_price": target,
		"active":       true,
		"triggered":    false,
	}).Decode(&existing)

	if err == nil {
		errs["_form"] = "You already have this exact active alert."
		return models.PriceAlert{}, errs
	}
	if err != nil && err != mongo.ErrNoDocuments {
		errs["_form"] = "Database error while creating the alert."
		return models.PriceAlert{}, errs
	}

	a := models.PriceAlert{
		UserID:      userID,
		Symbol:      sym,
		Condition:   cond,
		TargetPrice: target,

		Active:         true,
		Triggered:      false,
		TriggeredAt:    time.Time{},
		TriggeredPrice: 0,

		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	res, err := coll.InsertOne(ctx, a)
	if err != nil {
		errs["_form"] = "Could not create the alert."
		return models.PriceAlert{}, errs
	}

	a.ID = res.InsertedID.(primitive.ObjectID)
	return a, nil
}

func ListPriceAlerts(userID primitive.ObjectID, symbol string) ([]models.PriceAlert, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	sym := strings.ToUpper(strings.TrimSpace(symbol))
	coll := db.Client.Database("gomarket").Collection(alertsCollection)

	cur, err := coll.Find(ctx, bson.M{
		"user_id": userID,
		"symbol":  sym,
	}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]models.PriceAlert, 0)
	for cur.Next(ctx) {
		var a models.PriceAlert
		if err := cur.Decode(&a); err != nil {
			continue
		}
		out = append(out, a)
	}
	return out, nil
}

func DeletePriceAlert(userID primitive.ObjectID, alertID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	coll := db.Client.Database("gomarket").Collection(alertsCollection)
	_, err := coll.DeleteOne(ctx, bson.M{"_id": alertID, "user_id": userID})
	return err
}

// --- helpers for the background monitor ---

func ListActiveAlerts() ([]models.PriceAlert, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := db.Client.Database("gomarket").Collection(alertsCollection)

	cur, err := coll.Find(ctx, bson.M{
		"active":    true,
		"triggered": false,
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]models.PriceAlert, 0)
	for cur.Next(ctx) {
		var a models.PriceAlert
		if err := cur.Decode(&a); err != nil {
			continue
		}
		out = append(out, a)
	}
	return out, nil
}

func MarkAlertTriggered(alertID primitive.ObjectID, triggerPrice float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	coll := db.Client.Database("gomarket").Collection(alertsCollection)

	_, err := coll.UpdateOne(
		ctx,
		bson.M{"_id": alertID, "triggered": false},
		bson.M{"$set": bson.M{
			"active":          false,
			"triggered":       true,
			"triggered_at":    time.Now().UTC(),
			"triggered_price": roundToCents(triggerPrice),
			"updated_at":      time.Now().UTC(),
		}},
	)
	return err
}

func ListAllUserAlerts(userID primitive.ObjectID) ([]models.PriceAlert, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	coll := db.Client.Database("gomarket").Collection(alertsCollection)

	cur, err := coll.Find(ctx,
		bson.M{"user_id": userID},
		options.Find().SetSort(bson.D{
			{Key: "symbol", Value: 1},
			{Key: "created_at", Value: -1},
		}),
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]models.PriceAlert, 0)
	for cur.Next(ctx) {
		var a models.PriceAlert
		if err := cur.Decode(&a); err != nil {
			continue
		}
		out = append(out, a)
	}
	return out, nil
}
