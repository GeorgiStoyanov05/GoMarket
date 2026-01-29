package services

import (
	"errors"
	//"fmt"
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
	c.SetCookie("Auth", token, int(ttl), "/", "", false, true)
}

func ClearAuthCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", "", -1, "/", "", false, true)
}

func ChangeUserEmail(oldEmail string, newEmail string) (models.User, map[string]string){
	errs := map[string]string{}
    coll := db.Client.Database("gomarket").Collection("users")

    err := coll.FindOne(nil, bson.M{"email": newEmail}).Err()
	if err == nil {
		errs["email"] = "Email has already been taken!"
		return models.User{}, errs
	}

    _, err=coll.UpdateOne(nil, bson.M{"email": oldEmail}, bson.M{"$set":bson.M{"email": newEmail}})
    if err!=nil{
    	errs["_form"] = "There was a problem updating the email!"
     	return models.User{}, errs
    }

    var user models.User
        if err := coll.FindOne(nil, bson.M{"email": newEmail}).Decode(&user); err != nil {
            errs["_form"] = "Updated, but failed to load user."
            return models.User{}, errs
        }
    return user, nil
}

func ChangeUserPassword(id, password string) (models.User, map[string]string){
	errs := map[string]string{}
    coll := db.Client.Database("gomarket").Collection("users")
    var user models.User

    uid, err :=primitive.ObjectIDFromHex(id)
    if err!=nil{
    	errs["_form"] = "Error getting user Id!"
		return models.User{}, errs
    }
    err = coll.FindOne(nil, bson.M{"_id": uid}).Decode(&user)
	if err != nil {
		errs["_form"] = "Error getting the User!"
		return models.User{}, errs
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) == nil {
        errs["_form"] = "You've entered an old password!"
        return models.User{}, errs
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
       errs["_form"] = "Could not hash the new password."
       return models.User{}, errs
     }

    _, err=coll.UpdateOne(nil, bson.M{"_id": uid}, bson.M{"$set":bson.M{"password_hash": string(hash)}})
    if err!=nil{
    	errs["_form"] = "There was a problem updating the password!"
     	return models.User{}, errs
    }
    user.PasswordHash = string(hash)
    return user, nil
}

func ChangeUserBalance(id primitive.ObjectID, change float64) (models.User, map[string]string){
	errs := map[string]string{}
    coll := db.Client.Database("gomarket").Collection("users")

    _, err:=coll.UpdateOne(nil, bson.M{"_id": id}, bson.M{"$inc":bson.M{"balance": change}})
    if err!=nil{
    	errs["_form"] = "There was a problem updating the amount"
     	return models.User{}, errs
    }

    var user models.User
        if err := coll.FindOne(nil, bson.M{"_id": id}).Decode(&user); err != nil {
            errs["_form"] = "Updated, but failed to load user."
            return models.User{}, errs
        }
    return user, nil
}
