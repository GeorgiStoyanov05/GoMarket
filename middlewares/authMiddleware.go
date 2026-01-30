package middlewares

import (
	"errors"
	"net/http"
	"os"
	"time"
	db "github.com/GeorgiStoyanov05/GoMarket/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		v, exists := c.Get("IsLoggedIn")
		if !exists || v.(bool) == false {
			render404(c, "Unauthorized access!")
			return
		}
		c.Next()
	}
}

func isHTMX(c *gin.Context) bool {
	return c.GetHeader("HX-Request") == "true"
}

func render404(c *gin.Context, msg string) {
	data := WithAuth(c, gin.H{
		"message": msg,
		"View":    "404",
	})

	if v, ok := c.Get("IsLoggedIn"); ok { data["IsLoggedIn"] = v }
	if u, ok := c.Get("user"); ok { data["user"] = u }

	if isHTMX(c) {
		c.HTML(http.StatusOK, "404", data)
	} else {
		c.HTML(http.StatusNotFound, "index.html", data)
	}

	c.Abort()
}

func CheckIfLoggedIn() gin.HandlerFunc{
	return func(c *gin.Context) {

		c.Set("IsLoggedIn", false)
		tokenStr, err:=c.Cookie("Auth")
		if err!=nil{
			c.Next()
			return
		}

		token, err:=jwt.Parse(tokenStr, func(token *jwt.Token) (any, error){
			if _,ok:=token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unexpected signing method")
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err!=nil{
			c.Next()
			return
		}

		claims, ok:=token.Claims.(jwt.MapClaims)
		if !ok{
			c.Next()
			return
		}

		if claims["ttl"].(float64)<float64(time.Now().Unix()) {
			c.Next()
			return
		}
		userId:=claims["userID"].(string)
		transformedId,_ := primitive.ObjectIDFromHex(userId)
		user, ok:=db.GetUser(transformedId)

		if !ok{
			c.Next()
			return
		}

		c.Set("IsLoggedIn", true)
		c.Set("user", user)

		c.Next()
	}
}
