package controllers

import (
	"errors"
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"
)

//Login verifies cookie session for login
func Login(c *gin.Context) {
	// Get the email and password from request body
	var body models.InputBody

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Failed to read body"})
		return
	}

	// Look up requested user
	var user models.User
	database.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid email or password"})
		return
	}

	// Compare sent-in password with saved user password hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := generateJWTToken(user.ID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to generate token"})
		return
	}

	// Set token as a cookie
	c.SetCookie("token", token, 60*60, "/", "", false, true)

	// Return success response
	c.JSON(http.StatusOK, gin.H{"Message": "Login successful"})
}

//generateJWTToken generates a JWT token with a specified expiration time and user ID. It first sets the token expiration time based on the amountDays parameter passed into the function.
func generateJWTToken(userID uint, amountDays time.Duration) (string, error) {
	// Set token expiration time
	expirationTime := time.Now().Add(amountDays * 24 * time.Hour)

	// Create JWT claims
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   strconv.Itoa(int(userID)),
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
