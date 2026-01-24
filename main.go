package main

import (
	database "github.com/GeorgiStoyanov05/GoMarket2/database"
	//routes "github.com/GeorgiStoyanov05/GoMarket2/routes"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//routes.AuthRoutes(router)
	//routes.UserRoutes(router)

	database.DBInstance()

	router.Run(":" + port)
}
