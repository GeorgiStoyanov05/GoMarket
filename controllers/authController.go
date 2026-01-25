package controllers

import (
	"fmt"
	"net/http"

	models "github.com/GeorgiStoyanov05/GoMarket2/models"
	"github.com/gin-gonic/gin"
)

func GetRegisterPage(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(200, "register", gin.H{})
		return
	}
	c.HTML(200, "index.html", gin.H{
		"InitialPath": "/register",
	})
}

func PostRegisterPage(c *gin.Context) {
	var input models.RegisterModel
	errs := map[string]string{}
	if err := c.ShouldBind(&input); err != nil {
		errs["form"] = err.Error()

		c.HTML(http.StatusUnprocessableEntity, "register_form", gin.H{
			"values": input,
			"errors": errs,
		})
		return
	}
	c.Header("HX-Redirect", "/login")
	c.Status(http.StatusNoContent) // 204
}

func GetLoginPage(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(200, "login", gin.H{})
		return
	}
	c.HTML(200, "index.html", gin.H{
		"InitialPath": "/login",
	})
}

func PostLoginPage(c *gin.Context) {
}

func UserLogout(c *gin.Context) {
	fmt.Println("Logout")
}
