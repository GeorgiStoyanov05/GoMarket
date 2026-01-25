package models
type RegisterModel struct {
	First_Name        string	`form:"first_name" binding:"required"`
	Last_Name         string	`form:"last_name"  binding:"required"`
	Email             string	`form:"email" binding:"required,email"`
	Password          string	`form:"password" binding:"required,min=8"`
	RepeatPass		  string	`form:"rePassword" binding:"required"`
	RememberMe		  bool		`form:"rememberMe" binding:"required"`
}
