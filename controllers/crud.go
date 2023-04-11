package controllers

import (
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

//CreateName create new name on database of type NameType
func CreateName(c *gin.Context) {
	//name is passed by middlewares
	nameValue, ok := c.Get("name")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "name is not present on middlewares"})
		return
	}

	//parse nameValue into models.NameTypeInput
	var input models.NameType
	input, ok = nameValue.(models.NameType)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "failed to parse name"})
		return
	}

	//create name
	n, err := input.CreateName()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, n)
}

//GetID read name by id
func GetID(c *gin.Context) {
	var name models.NameType

	param := c.Params.ByName("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Id": "error on converting id"})
		return
	}

	n, _, err := name.GetNameById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}
	if n.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}

	c.JSON(http.StatusOK, n)
}

//GetName read name by name
func GetName(c *gin.Context) {
	var name models.NameType

	param := c.Params.ByName("name")
	n, err := name.GetNameByName(strings.ToUpper(param))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name not found"})
		return
	}
	if n.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name not found"})
		return
	}

	c.JSON(http.StatusOK, n)
	return
}

//UpdateName update name by id
func UpdateName(c *gin.Context) {
	//convert id string into int
	param := c.Params.ByName("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Id": "error on converting id"})
		return
	}

	//get the name by id
	var n models.NameType
	name, db, err := n.GetNameById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "name id not found"})
		return
	}
	if name.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}

	//name is passed by middlewares
	nameValue, ok := c.Get("name")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "name is not present on middlewares"})
		return
	}

	//parse nameValue into models.NameTypeInput
	var input models.NameType
	input, ok = nameValue.(models.NameType)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "failed to parse name"})
		return
	}

	if input.Name == name.Name && input.Classification == name.Classification && input.Metaphone == name.Metaphone && input.NameVariations == name.NameVariations {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Every item on json is the same on the database id " + param})
		return
	} else {
		if input.Name != "" {
			name.Name = input.Name
		}
		if input.Classification != "" {
			name.Classification = input.Classification
		}
		if input.Metaphone != "" {
			name.Metaphone = input.Metaphone
		}
		if input.NameVariations != "" {
			name.NameVariations = input.NameVariations
		}

		db.Save(name)
		c.JSON(http.StatusOK, name)
	}
}

//DeleteName delete name off database by id
func DeleteName(c *gin.Context) {
	var name models.NameType

	param := c.Params.ByName("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error on converting id"})
		return
	}

	n, _, err := name.GetNameById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}
	if n.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}

	_, err = name.DeleteNameById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Delete": "id " + param + " was deleted"})
}
