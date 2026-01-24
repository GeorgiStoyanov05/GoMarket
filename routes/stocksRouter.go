package routes

import (
	"github.com/gin-gonic/gin"
	"html/template"
)

func StocksRoutes(r *gin.Engine) {
	r.GET("/watchlist", func(c *gin.Context){
		tmpl:=template.Must(template.ParseFiles("views/index.html", "views/components/watchlist.html"))
		tmpl.Execute(c.Writer, nil)
	})
	r.GET("/portfolio", func(c *gin.Context){
		tmpl:=template.Must(template.ParseFiles("views/index.html", "views/components/portfolio.html"))
		tmpl.Execute(c.Writer, nil)
	})
}
