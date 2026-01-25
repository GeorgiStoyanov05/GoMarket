package routes

import (
	"github.com/GeorgiStoyanov05/GoMarket2/middlewares"
	"github.com/gin-gonic/gin"
)

func StocksRoutes(r *gin.Engine) {
	r.GET("/watchlist", func(c *gin.Context) {
		if c.GetHeader("HX-Request") == "true" {
			c.HTML(200, "watchlist", middlewares.WithAuth(c, gin.H{}))
			return
		}
		c.HTML(200, "index.html", middlewares.WithAuth(c, gin.H{
			"InitialPath": "/watchlist",
		}))
	})
	r.GET("/portfolio", func(c *gin.Context) {
		if c.GetHeader("HX-Request") == "true" {
			c.HTML(200, "portfolio", middlewares.WithAuth(c, gin.H{}))
			return
		}
		c.HTML(200, "index.html", middlewares.WithAuth(c, gin.H{
			"InitialPath": "/portfolio",
		}))
	})
}
