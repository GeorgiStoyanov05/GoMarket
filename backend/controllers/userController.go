package controllers

import (
	"net/http"
	"time"

	"backend/helpers"
	"backend/models"
	"backend/services"

	"github.com/gin-gonic/gin"
)

func UserSignUp(c *gin.Context) {
	var req models.SignUpModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.RegisterUser(c.Request.Context(), req)
	if err != nil {
		switch err {
		case services.ErrEmailTaken:
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	token, err := helpers.CreateAccessToken(user.ID.Hex(), user.Email, user.UserType, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	// HttpOnly cookie (recommended)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", token, 60*15, "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{"user": models.ToUserPublic(user)})
}

func UserSignIn(c *gin.Context) {
	var req models.SignInModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	token, err := helpers.CreateAccessToken(user.ID.Hex(), user.Email, user.UserType, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", token, 60*15, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"user": models.ToUserPublic(user)})
}

func UserLogout(c *gin.Context) {
	// clear cookie
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
