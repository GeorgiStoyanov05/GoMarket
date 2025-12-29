package routes

import (
	controller "backend/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controller.UserSignUp)
	}
}
