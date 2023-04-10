package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"time"
)

const MaxThreadsByToken = 5

// ValidateAuth returns a Gin middleware function that checks for a valid JWT token in the request header or cookie, and limits the rate of requests to`prevent DDoS attacks.
//	- The rate limit is enforced using a token bucket algorithm.
//	- The rate limit and queue capacity can be adjusted by modifying the constants in the function.
//	- If the token is invalid or has expired, or if the request cannot be processed due to an error, the middleware function aborts the request with a 401 Unauthorized HTTP status code.
func ValidateAuth() gin.HandlerFunc {
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
