package controllers

import (
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const levenshtein = 0.8

//CreateName create new name on database of type NameType
func CreateName(c *gin.Context) {
	var name models.NameType
	if err := c.ShouldBindJSON(&name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	database.Db.Create(&name)
	c.JSON(http.StatusOK, name)
}

//SearchNameByID read name by id
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

//DeleteName delete name off database by id
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

//UpdateName update name by id
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

//GetName read name by name
func GetName(c *gin.Context) {
	var name models.NameType

	n := c.Params.ByName("name")
	database.Db.Where("name = ?", strings.ToUpper(n)).Find(&name)

	if name.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name not found"})
		return
	}

	c.JSON(http.StatusOK, name)
}

//SearchSimilarNames search for all similar names by metaphone and Levenshtein method
func SearchSimilarNames(c *gin.Context) {
	var names []models.NameType

	//Name to be searched
	name := c.Params.ByName("name")
	database.Db.Find(&names)

	similarNames, mtf := findSimilarNames(names, name, levenshtein)

	//in case of failure in find a metaphone conde we return status not found
	if len(names) == 0 || len(similarNames) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "metaphone not found", "metaphone": mtf})
		return
	}

	//when the similar names result's in less than 5 we search for every similar name of all similar names founded previously
	if len(similarNames) < 5 {
		for _, n := range similarNames {
			similarNames, _ = findSimilarNames(names, n.Name, levenshtein)
		}
	}

	//order all similar names from high to low Levenshtein
	nameV := orderByLevenshtein(similarNames)

	//build canonical return
	canonicalEntity := findCanonical(name, nameV)
	r := models.MetaphoneR{
		ID:             canonicalEntity.ID,
		CreatedAt:      canonicalEntity.CreatedAt,
		UpdatedAt:      canonicalEntity.UpdatedAt,
		DeletedAt:      canonicalEntity.DeletedAt,
		Name:           canonicalEntity.Name,
		Classification: canonicalEntity.Classification,
		Metaphone:      canonicalEntity.Metaphone,
		NameVariations: nameV,
	}

	//return
	c.JSON(200, r)
}
