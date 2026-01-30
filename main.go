package main

import (
	"html/template"
	"os"
	"time"
	"context"
	"github.com/GeorgiStoyanov05/GoMarket2/services"
	database "github.com/GeorgiStoyanov05/GoMarket2/database"
	middlewares "github.com/GeorgiStoyanov05/GoMarket2/middlewares"
	routes "github.com/GeorgiStoyanov05/GoMarket2/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	router := gin.New()
	router.Static("/static", "./static")
	tmpl := template.Must(template.ParseGlob("views/*.html"))
	template.Must(tmpl.ParseGlob("views/components/*.html"))
	template.Must(tmpl.ParseGlob("views/components/partials/*.html"))
	router.SetHTMLTemplate(tmpl)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(middlewares.CheckIfLoggedIn())
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.HomeRoutes(router)
	routes.StocksRoutes(router)
	routes.AlertsRoutes(router)
	routes.TradingRoutes(router)
	router.NoRoute(func(c *gin.Context) {
		if c.GetHeader("HX-Request") == "true" {
			c.HTML(200, "404.html", gin.H{})
			return
		}
		c.HTML(200, "index.html", gin.H{
			"InitialPath": "/",
		})
	})
	database.DBInstance()
	services.StartPriceAlertMonitor(context.Background())
	services.EnsureTradingIndexes()
	router.Run(":" + port)
}
