package routes

import (
	"github.com/gin-gonic/gin"
)

func HomeRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		if c.GetHeader("HX-Request") == "true" {
			c.HTML(200, "home.html", gin.H{})
			return
		}
		c.HTML(200, "index.html", gin.H{
			"InitialPath": "/",
		})
	})
	r.GET("/search", func(c *gin.Context) {
		if c.GetHeader("HX-Request") == "true" {
			c.HTML(200, "search", gin.H{})
			return
		}
		c.HTML(200, "index.html", gin.H{
			"InitialPath": "/search",
		})
	})
}
