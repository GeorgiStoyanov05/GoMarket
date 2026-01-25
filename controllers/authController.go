package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
	"strings"
	models "github.com/GeorgiStoyanov05/GoMarket2/models"
	"github.com/GeorgiStoyanov05/GoMarket2/services"
	"github.com/gin-gonic/gin"
)

var emailRe = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+$`)

func isValidEmailRegex(s string) bool {
	s = strings.TrimSpace(s)
	return emailRe.MatchString(s)
}

func GetRegisterPage(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(200, "register", gin.H{
			"values": models.RegisterModel{},
			"errors": map[string]string{},
		})
		return
	}
	c.HTML(200, "index.html", gin.H{
		"InitialPath": "/register",
	})
}

func PostRegisterPage(c *gin.Context) {
	remember := c.PostForm("rememberMe") == "on"
	user := models.RegisterModel{
		First_Name: c.PostForm("first_name"),
		Last_Name:  c.PostForm("last_name"),
		Email:      c.PostForm("email"),
		Password:   c.PostForm("password"),
		RepeatPass: c.PostForm("rePassword"),
		RememberMe: remember,
	}

	errs := map[string]string{}

	if len(user.First_Name) < 2 {
		errs["first_name"] = "First name should be at least 2 characters long!"
	}
	if len(user.Last_Name) < 2 {
		errs["last_name"] = "Last name should be at least 2 characters long!"
	}
	if user.Email == "" {
		errs["email"] = "Email is required."
	}
	if !isValidEmailRegex(user.Email) {
		errs["email"] = "Please enter a valid email address."
	}
	if len(user.Password) < 6 {
		errs["password"] = "Password should be at least 6 characters long!"
	}
	if len(user.RepeatPass) < 6 {
		errs["rePassword"] = "Password should be at least 6 characters long!"
	}
	if user.Password != "" && user.RepeatPass != "" && user.Password != user.RepeatPass {
		errs["rePassword"] = "Passwords do not match."
	}

	if len(errs) > 0 {
		c.HTML(http.StatusOK, "register", gin.H{
			"values": user,
			"errors": errs,
		})
		return
	}

	u, newErrs := services.RegisterUser(&user)

	if len(newErrs) > 0 {
		c.HTML(http.StatusOK, "register", gin.H{
			"values": user,
			"errors": newErrs,
		})
		return
	}


	var ttl int64
	if(user.RememberMe==true){
		ttl=time.Now().Add(time.Hour*24*14).Unix()
	} else {
		ttl=time.Now().Add(time.Hour*2).Unix()
	}

	tokenStr, err:=services.CreateAndSignJWT(&u, ttl)
	if err!=nil{
		c.HTML(http.StatusOK, "login", gin.H{
			"values": user,
			"errors": map[string]string{"_form": err.Error()},
		})
		return
	}

	services.SetCookie(c, tokenStr, ttl)

	c.Header("HX-Redirect", "/")
	c.Status(204)
}

func GetLoginPage(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(200, "login", gin.H{
			"values": models.RegisterModel{},
			"errors": map[string]string{},
		})
		return
	}
	c.HTML(200, "index.html", gin.H{
		"InitialPath": "/login",
	})
}

func PostLoginPage(c *gin.Context) {
	remember := c.PostForm("rememberMe") == "on"
	user := models.LoginModel{
		Email:      c.PostForm("email"),
		Password:   c.PostForm("password"),
		RememberMe: remember,
	}
	errs := map[string]string{}

	if user.Email == "" {
		errs["email"] = "Email is required."
	}
	if !isValidEmailRegex(user.Email) {
		errs["email"] = "Please enter a valid email address."
	}

	if !isValidEmailRegex(user.Email) {
		errs["email"] = "Please enter a valid email address."
	}
	if len(user.Password) < 6 {
		errs["password"] = "Password should be at least 6 characters long!"
	}

	if len(errs) > 0 {
		c.HTML(http.StatusOK, "login", gin.H{
			"values": user,
			"errors": errs,
		})
		return
	}

	u,authErrs:=services.LoginUser(&user)
	if(len(authErrs)>0){
		c.HTML(http.StatusOK, "login", gin.H{
			"values": user,
			"errors": authErrs,
		})
		return
	}

	var ttl int64
	if(user.RememberMe==true){
		ttl=time.Now().Add(time.Hour*24*14).Unix()
	} else {
		ttl=time.Now().Add(time.Hour*2).Unix()
	}

	tokenStr, err:=services.CreateAndSignJWT(&u, ttl)
	if err!=nil{
		c.HTML(http.StatusOK, "login", gin.H{
			"values": user,
			"errors": map[string]string{"_form": err.Error()},
		})
		return
	}

	services.SetCookie(c, tokenStr, ttl)

	c.Header("HX-Redirect", "/")
	c.Status(204)
}

func UserLogout(c *gin.Context) {
	fmt.Println("Logout")
}
