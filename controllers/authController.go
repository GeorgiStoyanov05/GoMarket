package controllers

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

func GetRegisterPage(c *gin.Context){
tmpl:=template.Must(template.ParseFiles("views/index.html", "views/components/register.html"))
tmpl.Execute(c.Writer, nil)
}

func PostRegisterPage(c *gin.Context){
}

func GetLoginPage(c *gin.Context){
tmpl:=template.Must(template.ParseFiles("views/index.html", "views/components/login.html"))
tmpl.Execute(c.Writer, nil)
}

func PostLoginPage(c *gin.Context){
}

func UserLogout(c *gin.Context){
	fmt.Println("Logout")
}
