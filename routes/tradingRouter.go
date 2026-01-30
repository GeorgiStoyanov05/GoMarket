package routes

import (
	"github.com/GeorgiStoyanov05/GoMarket2/controllers"
	"github.com/GeorgiStoyanov05/GoMarket2/middlewares"
	"github.com/gin-gonic/gin"
)

func TradingRoutes(r *gin.Engine) {
	r.POST("/trade/:symbol/buy", middlewares.AuthMiddleware(), controllers.PostMarketBuy)
	r.GET("/positions/:symbol", middlewares.AuthMiddleware(), controllers.GetPositionPanel)
	r.POST("/trade/:symbol/sell", middlewares.AuthMiddleware(), controllers.PostMarketSell)
	r.GET("/portfolio", middlewares.AuthMiddleware(), controllers.GetPortfolioPage)
	r.GET("/portfolio/positions", middlewares.AuthMiddleware(), controllers.GetPortfolioPositions)

}
