package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const MaxThreadsByToken = 5

// RequireAuth returns a Gin middleware function that checks for a valid JWT token in the request header or cookie, and limits the rate of requests to`prevent DDoS attacks.
//	- The rate limit is enforced using a token bucket algorithm.
//	- The rate limit and queue capacity can be adjusted by modifying the constants in the function.
//	- If the token is invalid or has expired, or if the request cannot be processed due to an error, the middleware function aborts the request with a 401 Unauthorized HTTP status code.
func RequireAuth() gin.HandlerFunc {
	// Create a new rate limiter to limit the number of requests per second
	limiter := rate.NewLimiter(20000, MaxThreadsByToken)

	return func(c *gin.Context) {
		// Check if the request has exceeded the rate limit
		if !limiter.Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		// Get the token from the header or cookie
		tokenString := c.GetHeader("Token")
		var err error
		if tokenString == "" {
			tokenString, err = c.Cookie("token")
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		// Decode/validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check the expiration date
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// Continue
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

// ValidateIDParam validates id param. It must contain only numbers
func ValidateIDParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to parse the ":id" parameter as an integer
		if _, err := strconv.Atoi(c.Param("id")); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid ':id' parameter: '%s' is not a valid integer", c.Param("id")),
			})
			return
		}
		c.Next()
	}
}

//ValidateNameParam validates :name param. It must not contain numbers or spaces
func ValidateNameParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to retrieve the ":name" parameter from the request context
		name := c.Param("name")

		// Check if the name contains whitespace
		if strings.Contains(name, " ") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid ':name' parameter: '%s' should contain a single word with no spaces", name),
			})
			return
		}

		// Check if the name contains any numbers
		if _, err := strconv.Atoi(name); err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid ':name' parameter: '%s' should not contain any numbers", name),
			})
			return
		}

		c.Next()
	}
}
