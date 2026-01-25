package main

import (
	"html/template"
	"os"
	"time"

	database "github.com/GeorgiStoyanov05/GoMarket2/database"
	routes "github.com/GeorgiStoyanov05/GoMarket2/routes"
	middlewares "github.com/GeorgiStoyanov05/GoMarket2/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	tmpl := template.Must(template.ParseGlob("views/*.html"))
	template.Must(tmpl.ParseGlob("views/components/*.html"))
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

	router.Run(":" + port)
}
