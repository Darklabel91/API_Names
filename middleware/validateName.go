package middleware

import (
	"fmt"
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

//ValidateName validates :name param. It must not contain numbers or spaces
func ValidateName() gin.HandlerFunc {
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

//ValidateNameJSON validates JSON on models.NameType body
func ValidateNameJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		var name models.NameType
		if err := c.ShouldBindJSON(&name); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid JSON request body"})
			return
		}
		c.Set("name", name)
		c.Next()
	}
}
