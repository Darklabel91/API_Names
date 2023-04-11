package controllers

import (
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

//Signup a new user to the database
func Signup(c *gin.Context) {
	//get email/pass off req body
	var body models.UserInputBody

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Failed to read body"})
		return
	}

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Failed to hash password"})
		return
	}

	//create the user
	user := models.User{Email: body.Email, Password: string(hash), IP: c.ClientIP()}
	u, err := user.CreateUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Email already registered"})
		return
	}

	//respond
	c.JSON(http.StatusOK, gin.H{"Message": "User created", "User": u})
}
