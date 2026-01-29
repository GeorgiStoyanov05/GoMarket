package routes

import (
	"github.com/GeorgiStoyanov05/GoMarket2/controllers"
	"github.com/GeorgiStoyanov05/GoMarket2/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.GET("/settings", middlewares.AuthMiddleware(), controllers.GetUserSettings)
	r.GET("/settings/email", middlewares.AuthMiddleware(), controllers.GetChangeEmail)
	r.POST("/settings/email", middlewares.AuthMiddleware(), controllers.PostChangeEmail)
	r.GET("/settings/password", middlewares.AuthMiddleware(), controllers.GetChangePassword)
	r.POST("/settings/password", middlewares.AuthMiddleware(), controllers.PostChangePassword)
	r.GET("/funds", middlewares.AuthMiddleware(), controllers.GetFunds)
	r.POST("/funds", middlewares.AuthMiddleware(), controllers.PostFunds)
}
