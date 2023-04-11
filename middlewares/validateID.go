package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ValidateID validates id param. It must contain only numbers
func ValidateID() gin.HandlerFunc {
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
