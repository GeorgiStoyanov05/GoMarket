package routes

import (
	"github.com/GeorgiStoyanov05/GoMarket2/controllers"
	middlewares "github.com/GeorgiStoyanov05/GoMarket2/middlewares"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
		r.GET("/register", controllers.GetRegisterPage)
		r.POST("/register", controllers.PostRegisterPage)
		r.GET("/login", middlewares.AuthMiddleware, controllers.GetLoginPage)
		r.POST("/login", controllers.PostLoginPage)
		r.GET("/logout", controllers.UserLogout)
}
