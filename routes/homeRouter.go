package routes

import (
	"github.com/gin-gonic/gin"
	"html/template"
)

func HomeRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context){
		tmpl:=template.Must(template.ParseFiles("views/index.html", "views/components/home.html"))
		tmpl.Execute(c.Writer, nil)
	})
	r.GET("/search", func(c *gin.Context){
		tmpl:=template.Must(template.ParseFiles("views/index.html", "views/components/search.html"))
		tmpl.Execute(c.Writer, nil)
	})
}
