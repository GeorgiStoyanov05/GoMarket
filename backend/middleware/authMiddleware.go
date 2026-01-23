package middleware

import (
	"backend/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie("access_token")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		claims, err := helpers.ParseAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}
		c.Set("userID", claims.Subject)
		c.Set("email", claims.Email)
		c.Set("userType", claims.UserType)

		c.Next()
	}
}
