package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

const (
	TokenHeader = "Token"
	TokenCookie = "token"
)

const (
	MaxRequestsPerSecond = 5000
	MaxThreadsByToken    = 4
)

// ValidateAuth returns a Gin middleware function that checks for a valid JWT token in the request header or cookie, and aborts the request with a 401 Unauthorized HTTP status code if the token is invalid or has expired.
func ValidateAuth() gin.HandlerFunc {
	// Decode/validate the token
	return func(c *gin.Context) {
		// Get the token from the header or cookie
		tokenString := c.GetHeader(TokenHeader)
		if tokenString == "" {
			var err error
			tokenString, err = c.Cookie(TokenCookie)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

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

// RateLimit returns a Gin middleware function that limits the rate of requests to prevent DDoS attacks.
// The rate limit is enforced using a token bucket algorithm.
func RateLimit() gin.HandlerFunc {
	// Create a new rate limiter to limit the number of requests per second
	limiter := rate.NewLimiter(MaxRequestsPerSecond, MaxThreadsByToken)

	return func(c *gin.Context) {
		// Check if the request has exceeded the rate limit
		if !limiter.Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		// Continue
		c.Next()
	}
}
