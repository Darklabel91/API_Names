package controllers

import (
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllNames(c *gin.Context) {
	var names []models.NameType
	database.Db.Find(&names)

	c.JSON(200, names)
}

func CreateName(c *gin.Context) {
	var name models.NameType
	if err := c.ShouldBindJSON(&name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	database.Db.Create(&name)
	c.JSON(http.StatusOK, name)
}

func SearchNameByID(c *gin.Context) {
	var name models.NameType

	id := c.Params.ByName("id")
	database.Db.First(&name, id)

	if name.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}

	c.JSON(http.StatusOK, name)
}

func SearchNameByMetaphone(c *gin.Context) {
	var name models.NameType
	metaphone := c.Params.ByName("mtf")

	database.Db.Where(&models.NameType{Metaphone: metaphone}).First(&name)
	if name.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}

	c.JSON(http.StatusOK, name)
}

func DeleteName(c *gin.Context) {
	var name models.NameType
	id := c.Params.ByName("id")

	if name.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}

	database.Db.Delete(&name, id)
	c.JSON(http.StatusOK, gin.H{"data": "name data deleted"})
}

func UpdateName(c *gin.Context) {
	var name models.NameType
	id := c.Param("id")

	database.Db.First(&name, id)

	if err := c.ShouldBindJSON(&name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.Db.Model(&name).UpdateColumns(name)
	c.JSON(http.StatusOK, name)

}
