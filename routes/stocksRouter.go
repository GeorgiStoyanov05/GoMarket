package routes

import (
	"github.com/gin-gonic/gin"
)

func StocksRoutes(r *gin.Engine) {
	r.GET("/watchlist", func(c *gin.Context) {
		if c.GetHeader("HX-Request") == "true" {
			c.HTML(200, "watchlist", gin.H{})
			return
		}
		c.HTML(200, "index.html", gin.H{
			"InitialPath": "/watchlist",
		})
	})
	r.GET("/portfolio", func(c *gin.Context) {
		if c.GetHeader("HX-Request") == "true" {
			c.HTML(200, "portfolio", gin.H{})
			return
		}
		c.HTML(200, "index.html", gin.H{
			"InitialPath": "/portfolio",
		})
	})
}
