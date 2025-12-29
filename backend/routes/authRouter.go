package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.UserSignUp)
		auth.POST("/login", controllers.UserSignIn)
		auth.POST("/logout", controllers.UserLogout)

		// protected
		auth.GET("/me", middleware.RequireAuth(), controllers.Me)
	}
}
