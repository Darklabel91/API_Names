package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ValidateID is a Gin middleware function that validates the "id" parameter of the request URL. It checks that the parameter only contains valid integers.
func ValidateID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to parse the ":id" parameter as an integer
		if _, err := strconv.Atoi(c.Param("id")); err != nil {
			// If the parameter is not a valid integer, return a bad request error with a JSON response
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid ':id' parameter: '%s' is not a valid integer", c.Param("id")),
			})
			return
		}
		c.Next()
	}
}
