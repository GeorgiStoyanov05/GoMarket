// database/databaseConnect.go

package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GeorgiStoyanov05/GoMarket2/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Init() {
	Client = DBInstance()
}

func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MongoDB := os.Getenv("MONGODB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDB))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func GetUser(id primitive.ObjectID) (models.User, bool) {
	var u models.User
	coll := Client.Database("gomarket").Collection("users")
	err := coll.FindOne(nil, bson.M{"_id": id}).Decode(&u)
	if err != nil {
		return models.User{}, false
	}
	return u, true
}
