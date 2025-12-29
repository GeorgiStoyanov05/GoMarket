package middleware

import (
	"net/http"
	"strings"

	"backend/helpers"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) Try cookie
		token, _ := c.Cookie("access_token")

		// 2) Fallback: Authorization header (optional)
		if token == "" {
			auth := c.GetHeader("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				token = strings.TrimPrefix(auth, "Bearer ")
			}
		}

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

		// store info for handlers
		c.Set("userID", claims.Subject)
		c.Set("email", claims.Email)
		c.Set("userType", claims.UserType)

		c.Next()
	}
}
