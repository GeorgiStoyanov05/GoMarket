package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/GeorgiStoyanov05/GoMarket2/middlewares"
	"github.com/GeorgiStoyanov05/GoMarket2/models"
	"github.com/GeorgiStoyanov05/GoMarket2/services"
	"github.com/gin-gonic/gin"
)

func GetUserSettings(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(200, "settings", middlewares.WithAuth(c, gin.H{}))
		return
	}
	c.HTML(200, "index.html", middlewares.WithAuth(c, gin.H{
		"InitialPath": "/settings",
	}))
}

func GetChangeEmail(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(200, "changeEmail", middlewares.WithAuth(c, gin.H{
			"email":  "",
			"errors": map[string]string{},
		}))
		return
	}
	c.HTML(200, "index.html", middlewares.WithAuth(c, gin.H{
		"InitialPath": "/settings/email",
	}))
}

func PostChangeEmail(c *gin.Context) {
	newEmail := c.PostForm("email")
	errs := map[string]string{}
	uVal, ok := c.Get("user")
	if !ok {
		errs["_form"] = "There was an error getting user"
	}
	u, ok := uVal.(models.User)
	if !ok {
		errs["_form"] = "There was an error getting user"
	}
	if newEmail == "" {
		errs["email"] = "Email is required."
	}
	if !isValidEmailRegex(newEmail) {
		errs["email"] = "Please enter a valid email address."
	}

	if len(errs) > 0 {
		c.HTML(http.StatusOK, "changeEmail", middlewares.WithAuth(c, gin.H{
			"email":  newEmail,
			"errors": errs,
			"succ":   nil,
		}))
		return
	}

	u, newErrs := services.ChangeUserEmail(u.Email, newEmail)
	if len(newErrs) > 0 {
		c.HTML(http.StatusOK, "changeEmail", middlewares.WithAuth(c, gin.H{
			"email":  newEmail,
			"errors": newErrs,
			"succ":   nil,
		}))
		return
	}
	c.Set("user", u)
	c.HTML(http.StatusOK, "changeEmail", middlewares.WithAuth(c, gin.H{
		"email":  "",
		"errors": newErrs,
		"succ":   "You have changed your email successfully!",
	}))
}

func GetChangePassword(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(200, "changePassword", middlewares.WithAuth(c, gin.H{
			"errors": map[string]string{},
			"succ":   "",
		}))
		return
	}
	c.HTML(200, "index.html", middlewares.WithAuth(c, gin.H{
		"InitialPath": "/settings/password",
	}))
}

func PostChangePassword(c *gin.Context) {
	password := c.PostForm("password")
	rePassword := c.PostForm("rePassword")
	errs := map[string]string{}
	uVal, ok := c.Get("user")
	if !ok {
		errs["_form"] = "There was an error getting user"
	}
	user, ok := uVal.(models.User)
	if !ok {
		errs["_form"] = "There was an error getting user"
	}
	if len(password) < 6 {
		errs["password"] = "Password should be at least 6 characters long!"
	}
	if password != rePassword {
		errs["rePassword"] = "Passwords do not match!"
	}
	if len(errs) > 0 {
		c.HTML(200, "changePassword", middlewares.WithAuth(c, gin.H{
			"errors": errs,
			"succ":   "",
		}))
		return
	}
	u, newErrs := services.ChangeUserPassword(user.ID.Hex(), password)

	if len(newErrs) > 0 {
		c.HTML(http.StatusOK, "changePassword", middlewares.WithAuth(c, gin.H{
			"errors": newErrs,
			"succ":   "",
		}))
		return
	}

	c.Set("user", u)
	c.HTML(http.StatusOK, "changePassword", middlewares.WithAuth(c, gin.H{
		"errors": newErrs,
		"succ":   "You have changed your password successfully!",
	}))
}

func GetFunds(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(200, "depositFunds", middlewares.WithAuth(c, gin.H{
			"errors": map[string]string{},
			"amount": 0,
			"succ": "",
		}))
		return
	}
	c.HTML(200, "index.html", middlewares.WithAuth(c, gin.H{
		"InitialPath": "/funds",
	}))
}

func PostFunds(c *gin.Context) {
	errs := map[string]string{}
	amountStr := strings.TrimSpace(c.PostForm("amount"))
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err!=nil{
		errs["_form"] = "There was an error with the amount!"
	}
	uVal, ok := c.Get("user")
	if !ok {
		errs["_form"] = "Error getting user!"
	}
	user, ok := uVal.(models.User)
	if !ok {
		errs["_form"] = "There was an error getting user"
	}
	if amount <= 0 {
		errs["amount"] = "Amount must be bigger than zero!"
	}

	if len(errs) > 0 {
		c.HTML(http.StatusOK, "depositFunds", middlewares.WithAuth(c, gin.H{
			"errors": errs,
			"amount": amount,
			"succ":   "",
		}))
		return
	}

	u, newErrs := services.ChangeUserBalance(user.ID, amount)

	if len(newErrs) > 0 {
		c.HTML(http.StatusOK, "depositFunds", middlewares.WithAuth(c, gin.H{
			"errors": newErrs,
			"succ":   "",
			"amount": amount,
		}))
		return
	}

	c.Set("user", u)
	c.HTML(http.StatusOK, "depositFunds", middlewares.WithAuth(c, gin.H{
		"errors": newErrs,
		"succ":   "The deposit was successful!",
		"amount": 0,
	}))
}
