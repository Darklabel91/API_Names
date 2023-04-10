package controllers

import (
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

//GetName read name by name
func GetName(c *gin.Context) {
	var name models.NameType

	n := c.Params.ByName("name")
	database.DB.Raw("select * from name_types where name = ?", strings.ToUpper(n)).Find(&name)

	if name.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name not found"})
		return
	}

	c.JSON(http.StatusOK, name)
	return
}
