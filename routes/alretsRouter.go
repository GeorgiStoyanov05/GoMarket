package routes

import (
	"github.com/GeorgiStoyanov05/GoMarket2/controllers"
	"github.com/GeorgiStoyanov05/GoMarket2/middlewares"
	"github.com/gin-gonic/gin"
)

func AlertsRoutes(r *gin.Engine) {
	r.POST("/alerts/:symbol", middlewares.AuthMiddleware(), controllers.PostCreateAlert)
	r.GET("/alerts/:symbol/list", middlewares.AuthMiddleware(), controllers.GetAlertsList)
	r.POST("/alerts/:symbol/:id/delete", middlewares.AuthMiddleware(), controllers.PostDeleteAlert)
	r.GET("/alerts/list", middlewares.AuthMiddleware(), controllers.GetWatchlistAlerts)
	r.POST("/alerts/by-id/:id/delete", middlewares.AuthMiddleware(), controllers.PostDeleteAlertGlobal)
}
