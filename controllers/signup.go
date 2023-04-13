package controllers

import (
	"net/http"

	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Signup creates a new user and saves it to the database.
func Signup(c *gin.Context) {
	// Get email and password from request body.
	var body models.UserInputBody
	if c.Bind(&body) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
		return
	}

	// Hash the password.
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error hashing password"})
		return
	}

	// Create the user with hashed password and IP address of client.
	user := models.User{
		Email:    body.Email,
		Password: string(hash),
		IP:       c.ClientIP(),
	}
	u, err := user.CreateUser()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Respond with success message and created user.
	c.JSON(http.StatusOK, gin.H{"Message": "User created", "User": u})
}
