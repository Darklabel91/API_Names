package controllers

import (
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

//DeleteName delete name off database by id
func DeleteName(c *gin.Context) {
	var name models.NameType

	id := c.Params.ByName("id")
	database.DB.First(&name, id)

	if name.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}

	database.DB.Delete(&name, id)
	c.JSON(http.StatusOK, gin.H{"Delete": "name id " + id + " was deleted"})
}
