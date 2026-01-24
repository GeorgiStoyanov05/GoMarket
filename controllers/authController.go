package controllers

import (
	"fmt"
	"html/template"
	models "github.com/GeorgiStoyanov05/GoMarket2/models"
	"github.com/gin-gonic/gin"
)

func GetRegisterPage(c *gin.Context){
tmpl:=template.Must(template.ParseFiles("views/index.html", "views/components/register.html"))
tmpl.Execute(c.Writer, nil)
}

func PostRegisterPage(c *gin.Context){
	var input models.RegisterModel
	if err:=c.ShouldBind(&input); err!=nil{
	 c.JSON(400, gin.H{"error": err.Error()})
        return
	}
	fmt.Println(input.RememberMe)
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
