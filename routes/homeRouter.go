package routes

import (
	"github.com/GeorgiStoyanov05/GoMarket2/middlewares"
	"github.com/gin-gonic/gin"
)

func HomeRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		if c.GetHeader("HX-Request") == "true" {
			c.HTML(200, "home", middlewares.WithAuth(c, gin.H{}))
			return
		}
		c.HTML(200, "index.html", middlewares.WithAuth(c, gin.H{
			"InitialPath": "/",
		}))
	})
}
