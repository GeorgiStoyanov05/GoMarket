package main

import (
	routes "backend/routes"
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
	router.Use(gin.Logger())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true, // keep true if you will use cookies for JWT
		MaxAge:           12 * time.Hour,
	}))

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(c *gin.Context){
		c.JSON(200, gin.H{"success":"Access granted for api-1"})

		router.GET("/api-2", func(c *gin.Context){
			c.JSON(200, gin.H{"success": "Access granted for api-2"})
		})
	})

	router.Run(":"+port)
}
