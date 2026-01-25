package services

import (
	"errors"
	"net/http"
	"os"
	"time"

	db "github.com/GeorgiStoyanov05/GoMarket2/database"
	models "github.com/GeorgiStoyanov05/GoMarket2/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(user *models.RegisterModel) (models.User, map[string]string) {
	errs := map[string]string{}
	coll := db.Client.Database("gomarket").Collection("users")
	err := coll.FindOne(nil, bson.M{"email": user.Email}).Err()
	if err == nil {
		errs["email"] = "Email has already been taken!"
		return models.User{}, errs
	}
	if err != mongo.ErrNoDocuments {
		errs["_form"] = "There is a problem registering this user!"
		return models.User{}, errs
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, errs
	}

	var role string
	if user.Email == os.Getenv("ADMIN_EMAIL") && user.Password == os.Getenv("ADMIN_PASSWORD") {
		role = "Admin"
	} else {
		role = "User"
	}

	u := models.User{
		FirstName:    user.First_Name,
		LastName:     user.Last_Name,
		Email:        user.Email,
		PasswordHash: string(hash),
		Role:         role,
		Balance:      0,

		Watchlist: []models.Watchlist{},
		Portfolio: []models.BoughtStock{},

		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	res, err := coll.InsertOne(nil, u)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			errs["_form"] = "There was a problem creating the user!"
		}
		return models.User{}, errs
	}
	u.ID = res.InsertedID.(primitive.ObjectID)
	return u, nil
}


func LoginUser(user *models.LoginModel) (models.User, map[string]string){
 errs := map[string]string{}

    coll := db.Client.Database("gomarket").Collection("users")

    var u models.User
    err := coll.FindOne(nil, bson.M{"email": user.Email}).Decode(&u)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            errs["_form"] = "Invalid email or password."
            return models.User{}, errs
        }
        errs["_form"] = "Server error. Please try again."
        return models.User{}, errs
    }

    if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(user.Password)) != nil {
        errs["_form"] = "Invalid email or password."
        return models.User{}, errs
    }

    return u, nil
}

func CreateAndSignJWT(user *models.User, ttl int64) (string,error){
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID.Hex(),
		"ttl":	ttl,
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func SetCookie(c *gin.Context, token string, ttl int64){
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", token, int(ttl), "", "", false, true)
}
