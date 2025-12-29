package services

import (
	"context"
	"errors"
	"time"
	"strings"
	"backend/database"
	"backend/models"

	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrEmailTaken          = errors.New("email already registered")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)

func RegisterUser(ctx context.Context, req models.SignUpModel) (models.User, error) {
	coll := database.Client.Database("gomarket").Collection("users")
	err := coll.FindOne(ctx, bson.M{"email": req.Email}).Err()
	if err == nil {
		return models.User{}, ErrEmailTaken
	}
	if err != mongo.ErrNoDocuments {
		return models.User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	now := time.Now().UTC()

	user := models.User{
		ID:                primitive.NewObjectID(),
		FullName:          req.FullName,
		Email:             req.Email,
		PasswordHash:      string(hash),
		Country:           req.Country,
		PreferredIndustry: req.PreferredIndustry,
		UserType:          "USER",
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return models.User{}, ErrEmailTaken
		}
		return models.User{}, err
	}

	return user, nil
}

func LoginUser(ctx context.Context, email string, password string) (models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	coll := database.Client.Database("gomarket").Collection("users")

	email = strings.ToLower(strings.TrimSpace(email))

	var user models.User
	err := coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, ErrInvalidCredentials
		}
		return models.User{}, err
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return models.User{}, ErrInvalidCredentials
	}

	return user, nil
}
