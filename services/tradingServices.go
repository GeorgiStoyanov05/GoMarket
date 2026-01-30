package services

import (
	"context"
	"math"
	"strings"
	"time"

	db "github.com/GeorgiStoyanov05/GoMarket2/database"
	"github.com/GeorgiStoyanov05/GoMarket2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BuyResult struct {
	Symbol     string
	Qty        int64
	FillPrice  float64
	Cost       float64
	NewBalance float64
	Position   models.Position
}

type SellResult struct {
	Symbol     string
	Qty        int64
	FillPrice  float64
	Proceeds   float64
	NewBalance float64
	Remaining  *models.Position // nil if position closed
}

func roundMoney(v float64) float64 {
	return math.Round(v*100) / 100
}

func MarketBuy(userID primitive.ObjectID, symbol string, qty int64) (BuyResult, map[string]string) {
	errs := map[string]string{}

	sym := strings.ToUpper(strings.TrimSpace(symbol))
	if sym == "" {
		errs["symbol"] = "Missing symbol."
	}
	if qty <= 0 {
		errs["qty"] = "Quantity must be greater than 0."
	}
	if len(errs) > 0 {
		return BuyResult{}, errs
	}

	// 1) Get fill price (market buy = quote at click time)
	price, err := FetchCurrentPrice(sym)
	if err != nil {
		errs["_form"] = "Could not fetch current price."
		return BuyResult{}, errs
	}
	price = roundMoney(price)

	cost := roundMoney(price * float64(qty))

	// Try transaction (works on Atlas/replica set). If not supported, fallback to sequential.
	// Try transaction (works on Atlas/replica set). If not supported, fallback.
	if ok := tryMarketBuyTxn(userID, sym, qty, price, cost); ok != nil {
		if ok.NewBalance < 0 {
			return BuyResult{}, map[string]string{"balance": "Not enough balance for this purchase."}
		}
		return *ok, nil
	}

	// Fallback
	return marketBuyNoTxn(userID, sym, qty, price, cost)
}

func tryMarketBuyTxn(userID primitive.ObjectID, sym string, qty int64, price, cost float64) *BuyResult {
	client := db.Client
	sess, err := client.StartSession()
	if err != nil {
		return nil
	}
	defer sess.EndSession(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersColl := client.Database("gomarket").Collection("users")
	posColl := client.Database("gomarket").Collection("positions")
	ordersColl := client.Database("gomarket").Collection("orders")

	now := time.Now().UTC()

	var out BuyResult
	_, txnErr := sess.WithTransaction(ctx, func(sc mongo.SessionContext) (any, error) {
		// A) Deduct balance atomically only if enough money
		filter := bson.M{"_id": userID, "balance": bson.M{"$gte": cost}}
		update := bson.M{
			"$inc": bson.M{"balance": -cost},
			"$set": bson.M{"updated_at": now},
		}

		var updatedUser models.User
		err := usersColl.FindOneAndUpdate(
			sc,
			filter,
			update,
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		).Decode(&updatedUser)

		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		if err != nil {
			return nil, err
		}

		// B) Upsert position with atomic avg-cost math (pipeline update)
		var updatedPos models.Position
		err = posColl.FindOneAndUpdate(
			sc,
			bson.M{"user_id": userID, "symbol": sym},
			mongo.Pipeline{
				{{Key: "$set", Value: bson.D{
					{Key: "user_id", Value: userID},
					{Key: "symbol", Value: sym},
					{Key: "created_at", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$created_at", now}}}},
					{Key: "updated_at", Value: now},
					// qty and avg_cost computed from OLD values
					{Key: "qty", Value: bson.D{{Key: "$add", Value: bson.A{
						bson.D{{Key: "$ifNull", Value: bson.A{"$qty", 0}}},
						qty,
					}}}},
					{Key: "avg_cost", Value: bson.D{{Key: "$let", Value: bson.D{
						{Key: "vars", Value: bson.D{
							{Key: "oldQty", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$qty", 0}}}},
							{Key: "oldAvg", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$avg_cost", 0}}}},
							{Key: "buyQty", Value: qty},
							{Key: "buyPrice", Value: price},
						}},
						{Key: "in", Value: bson.D{{Key: "$cond", Value: bson.A{
							bson.D{{Key: "$gt", Value: bson.A{
								bson.D{{Key: "$add", Value: bson.A{"$$oldQty", "$$buyQty"}}}, 0,
							}}},
							bson.D{{Key: "$divide", Value: bson.A{
								bson.D{{Key: "$add", Value: bson.A{
									bson.D{{Key: "$multiply", Value: bson.A{"$$oldQty", "$$oldAvg"}}},
									bson.D{{Key: "$multiply", Value: bson.A{"$$buyQty", "$$buyPrice"}}},
								}}},
								bson.D{{Key: "$add", Value: bson.A{"$$oldQty", "$$buyQty"}}},
							}}},
							0,
						}}}},
					}}}},
				}}},
			},
			options.FindOneAndUpdate().
				SetUpsert(true).
				SetReturnDocument(options.After),
		).Decode(&updatedPos)

		if err != nil {
			return nil, err
		}

		// C) Insert order (ledger)
		order := models.Order{
			UserID:    userID,
			Symbol:    sym,
			Side:      "buy",
			Qty:       qty,
			Price:     price,
			CreatedAt: now,
		}
		if _, err := ordersColl.InsertOne(sc, order); err != nil {
			return nil, err
		}

		out = BuyResult{
			Symbol:     sym,
			Qty:        qty,
			FillPrice:  price,
			Cost:       cost,
			NewBalance: roundMoney(updatedUser.Balance),
			Position:   updatedPos,
		}
		return nil, nil
	})

	// If txn unsupported (common on local standalone), return nil to fallback
	if txnErr != nil {
		// insufficient funds case
		if txnErr == mongo.ErrNoDocuments {
			return &BuyResult{Symbol: sym, Qty: qty, FillPrice: price, Cost: cost, NewBalance: -1}
		}
		// treat any txn error as "fallback"
		return nil
	}

	return &out
}

func marketBuyNoTxn(userID primitive.ObjectID, sym string, qty int64, price, cost float64) (BuyResult, map[string]string) {
	errs := map[string]string{}
	client := db.Client

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersColl := client.Database("gomarket").Collection("users")
	posColl := client.Database("gomarket").Collection("positions")
	ordersColl := client.Database("gomarket").Collection("orders")

	now := time.Now().UTC()

	// A) Deduct balance if enough
	var updatedUser models.User
	err := usersColl.FindOneAndUpdate(
		ctx,
		bson.M{"_id": userID, "balance": bson.M{"$gte": cost}},
		bson.M{"$inc": bson.M{"balance": -cost}, "$set": bson.M{"updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedUser)

	if err == mongo.ErrNoDocuments {
		errs["balance"] = "Not enough balance for this purchase."
		return BuyResult{}, errs
	}
	if err != nil {
		errs["_form"] = "Database error while updating balance."
		return BuyResult{}, errs
	}

	// B) Upsert position (same pipeline)
	var updatedPos models.Position
	err = posColl.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": userID, "symbol": sym},
		mongo.Pipeline{
			{{Key: "$set", Value: bson.D{
				{Key: "user_id", Value: userID},
				{Key: "symbol", Value: sym},
				{Key: "created_at", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$created_at", now}}}},
				{Key: "updated_at", Value: now},
				{Key: "qty", Value: bson.D{{Key: "$add", Value: bson.A{
					bson.D{{Key: "$ifNull", Value: bson.A{"$qty", 0}}},
					qty,
				}}}},
				{Key: "avg_cost", Value: bson.D{{Key: "$let", Value: bson.D{
					{Key: "vars", Value: bson.D{
						{Key: "oldQty", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$qty", 0}}}},
						{Key: "oldAvg", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$avg_cost", 0}}}},
						{Key: "buyQty", Value: qty},
						{Key: "buyPrice", Value: price},
					}},
					{Key: "in", Value: bson.D{{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$gt", Value: bson.A{
							bson.D{{Key: "$add", Value: bson.A{"$$oldQty", "$$buyQty"}}}, 0,
						}}},
						bson.D{{Key: "$divide", Value: bson.A{
							bson.D{{Key: "$add", Value: bson.A{
								bson.D{{Key: "$multiply", Value: bson.A{"$$oldQty", "$$oldAvg"}}},
								bson.D{{Key: "$multiply", Value: bson.A{"$$buyQty", "$$buyPrice"}}},
							}}},
							bson.D{{Key: "$add", Value: bson.A{"$$oldQty", "$$buyQty"}}},
						}}},
						0,
					}}}},
				}}}},
			}}},
		},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&updatedPos)

	if err != nil {
		// Worst-case inconsistency: balance deducted but position failed.
		// For MVP, we surface an error; later we’ll add txn-only or rollback logic.
		errs["_form"] = "Bought balance updated, but position update failed."
		return BuyResult{}, errs
	}

	// C) Insert order
	order := models.Order{
		UserID:    userID,
		Symbol:    sym,
		Side:      "buy",
		Qty:       qty,
		Price:     price,
		CreatedAt: now,
	}
	if _, err := ordersColl.InsertOne(ctx, order); err != nil {
		errs["_form"] = "Purchase saved partially (order insert failed)."
		return BuyResult{}, errs
	}

	return BuyResult{
		Symbol:     sym,
		Qty:        qty,
		FillPrice:  price,
		Cost:       cost,
		NewBalance: roundMoney(updatedUser.Balance),
		Position:   updatedPos,
	}, nil
}

func MarketSell(userID primitive.ObjectID, symbol string, qty int64) (SellResult, map[string]string) {
	errs := map[string]string{}

	sym := strings.ToUpper(strings.TrimSpace(symbol))
	if sym == "" {
		errs["symbol"] = "Missing symbol."
	}
	if qty <= 0 {
		errs["qty"] = "Quantity must be greater than 0."
	}
	if len(errs) > 0 {
		return SellResult{}, errs
	}

	// Fill price from quote
	price, err := FetchCurrentPrice(sym)
	if err != nil {
		errs["_form"] = "Could not fetch current price."
		return SellResult{}, errs
	}
	price = roundMoney(price)

	proceeds := roundMoney(price * float64(qty))

	// For now: no transactions to keep this chunk smaller.
	// We'll do a safe sequential flow with validation using positions.
	return marketSellNoTxn(userID, sym, qty, price, proceeds)
}

func marketSellNoTxn(userID primitive.ObjectID, sym string, qty int64, price, proceeds float64) (SellResult, map[string]string) {
	errs := map[string]string{}
	client := db.Client

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersColl := client.Database("gomarket").Collection("users")
	posColl := client.Database("gomarket").Collection("positions")
	ordersColl := client.Database("gomarket").Collection("orders")

	// 1) Ensure user has enough shares (atomic-ish: check and decrement with filter)
	now := time.Now().UTC()

	// Decrement qty if enough
	updateRes := posColl.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": userID, "symbol": sym, "qty": bson.M{"$gte": qty}},
		bson.M{
			"$inc": bson.M{"qty": -qty},
			"$set": bson.M{"updated_at": now},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedPos models.Position
	err := updateRes.Decode(&updatedPos)
	if err == mongo.ErrNoDocuments {
		errs["qty"] = "You don't have enough shares to sell."
		return SellResult{}, errs
	}
	if err != nil {
		errs["_form"] = "Database error while selling."
		return SellResult{}, errs
	}

	// If qty hit 0, delete the position doc
	var remaining *models.Position
	if updatedPos.Qty <= 0 {
		_, _ = posColl.DeleteOne(ctx, bson.M{"_id": updatedPos.ID, "user_id": userID})
		remaining = nil
	} else {
		remaining = &updatedPos
	}

	// 2) Credit balance
	var updatedUser models.User
	err = usersColl.FindOneAndUpdate(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$inc": bson.M{"balance": proceeds}, "$set": bson.M{"updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedUser)
	if err != nil {
		// Worst-case inconsistency: position decreased but balance update failed.
		// We'll surface error; later we can wrap in txn for full safety.
		errs["_form"] = "Sold shares, but balance update failed."
		return SellResult{}, errs
	}

	// 3) Insert order
	order := models.Order{
		UserID:    userID,
		Symbol:    sym,
		Side:      "sell",
		Qty:       qty,
		Price:     price,
		CreatedAt: now,
	}
	if _, err := ordersColl.InsertOne(ctx, order); err != nil {
		errs["_form"] = "Sold, but failed to record order."
		return SellResult{}, errs
	}

	return SellResult{
		Symbol:     sym,
		Qty:        qty,
		FillPrice:  price,
		Proceeds:   proceeds,
		NewBalance: roundMoney(updatedUser.Balance),
		Remaining:  remaining,
	}, nil
}
