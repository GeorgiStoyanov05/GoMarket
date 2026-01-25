package main

import (
	"html/template"
	"os"
	"time"

	database "github.com/GeorgiStoyanov05/GoMarket2/database"
	routes "github.com/GeorgiStoyanov05/GoMarket2/routes"
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

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.HomeRoutes(router)
	routes.StocksRoutes(router)
	router.NoRoute(func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles("views/index.html", "views/components/404.html"))
		tmpl.Execute(c.Writer, nil)
	})
	database.DBInstance()

	router.Run(":" + port)
}
