package middlewares

import (
	"errors"
	"net/http"
	"os"
	"time"

	db "github.com/GeorgiStoyanov05/GoMarket2/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware(c *gin.Context){

	tokenStr, err:=c.Cookie("Auth")
	if err!=nil{
		render404(c, "Unauthorized access!")
		return
	}

	token, err:=jwt.Parse(tokenStr, func(token *jwt.Token) (any, error){
		if _,ok:=token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err!=nil{
		render404(c, err.Error())
		return
	}

	claims, ok:=token.Claims.(jwt.MapClaims)
	if !ok{
		render404(c, "JWT claims failed!")
		return
	}

	if claims["ttl"].(float64)<float64(time.Now().Unix()) {
		render404(c, "Token expired!")
		return
	}
	userId:=claims["userID"].(string)
	transformedId,_ := primitive.ObjectIDFromHex(userId)
	user, ok:=db.GetUser(transformedId)

	if !ok{
		render404(c, "User not found!")
		return
	}

	c.Set("user", user)

	c.Next()
}

func isHTMX(c *gin.Context) bool {
	return c.GetHeader("HX-Request") == "true"
}

func render404(c *gin.Context, msg string) {
	data := gin.H{
		"message": msg,
		"View":    "404",
	}

	if isHTMX(c) {
		c.HTML(http.StatusOK, "404", data)
	} else {
		c.HTML(http.StatusNotFound, "index.html", data)
	}

	c.Abort()
}
