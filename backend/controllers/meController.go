package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"userID":   c.GetString("userID"),
		"email":    c.GetString("email"),
		"userType": c.GetString("userType"),
	})
}
