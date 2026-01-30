package routes

import (
	"github.com/GeorgiStoyanov05/GoMarket2/middlewares"
	"github.com/GeorgiStoyanov05/GoMarket2/controllers"
	"github.com/gin-gonic/gin"
)

func StocksRoutes(r *gin.Engine) {
	r.GET("/watchlist", middlewares.AuthMiddleware(), func(c *gin.Context) {
		if c.GetHeader("HX-Request") == "true" {
			c.HTML(200, "watchlist", middlewares.WithAuth(c, gin.H{}))
			return
		}
		c.HTML(200, "index.html", middlewares.WithAuth(c, gin.H{
			"InitialPath": "/watchlist",
		}))
	})
	r.GET("/search", middlewares.AuthMiddleware(), func(c *gin.Context) {
		if c.GetHeader("HX-Request") == "true" {
			c.HTML(200, "search", middlewares.WithAuth(c, gin.H{}))
			return
		}
		c.HTML(200, "index.html",middlewares.WithAuth(c, gin.H{
			"InitialPath": "/search",
		}))
	})
	r.GET("/details/:symbol", middlewares.AuthMiddleware(), controllers.GetSymbolDetailsPage)
	r.GET("/ws/trades", controllers.WSFinnhubTrades)


}
