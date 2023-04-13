package controllers

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Login verifies email and password and sets a JWT token as a cookie for authentication
func Login(c *gin.Context) {
	// Get email and password from request body
	var body models.UserInputBody
	if err := c.Bind(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
		return
	}

	// Look up user by email
	var user models.User
	u, err := user.GetUserByEmail(body.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	// Compare password from request body with user's hashed password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(body.Password))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	token, err := generateJWTToken(u.ID, 1*time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set token as a cookie
	c.SetCookie("token", token, int(1*time.Hour.Seconds()), "/", "", false, true)

	// Return success response
	c.JSON(http.StatusOK, gin.H{"Message": "Login successful"})
}

//generateJWTToken generates a JWT token with a specified expiration time and user ID. It first sets the token expiration time based on the amountDays parameter passed into the function.
func generateJWTToken(userID uint, amountDays time.Duration) (string, error) {
	// Set token expiration time
	expirationTime := time.Now().Add(amountDays * 24 * time.Hour)

	// Create JWT claims
	claims := jwt.MapClaims{
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
		"sub": strconv.Itoa(int(userID)),
	}

	// Create token using claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token using secret key
	secretKey := []byte(os.Getenv("SECRET"))
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.New("failed to sign token")
	}

	return signedToken, nil
}
