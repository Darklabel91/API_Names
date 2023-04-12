package controllers

import (
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

//CreateName creates a new name in the database of type NameType
func CreateName(c *gin.Context) {
	// The name is passed by middlewares
	nameValue, ok := c.Get("name")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "name is not present on middlewares"})
		return
	}

	// Parse nameValue into models.NameTypeInput
	var input models.NameType
	input, ok = nameValue.(models.NameType)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "failed to parse name"})
		return
	}

	// Check the cache
	preloadTable := checkCache(c)

	// Check if there's an exact name on the database
	for _, name := range preloadTable {
		if name.Name == input.Name {
			c.JSON(http.StatusBadRequest, gin.H{"Message": " Duplicate entry " + name.Name})
			return
		}
	}

	// Create name
	n, err := input.CreateName()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	// Clear cache
	clearCache(c)

	// Return successful response
	c.JSON(http.StatusOK, n)
	return
}

//GetID reads a name by id
func GetID(c *gin.Context) {
	var name models.NameType

	// Convert id string into int
	param := c.Params.ByName("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Id": "error on converting id"})
		return
	}

	// Get the name by id
	n, _, err := name.GetNameById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}
	if n.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, n)
	return
}

//GetName reads a name by name
func GetName(c *gin.Context) {
	var name models.NameType

	// Get name to be searched
	param := c.Params.ByName("name")

	// Search for name
	n, err := name.GetNameByName(strings.ToUpper(param))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name not found"})
		return
	}
	if n.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name not found"})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, n)
	return
}

//GetMetaphoneMatch reads a name by metaphone
func GetMetaphoneMatch(c *gin.Context) {
	var nameType models.NameType

	// Get name to be searched
	name := c.Params.ByName("name")

	// Check the cache
	preloadTable := checkCache(c)

	// Search for similar names
	canonicalEntity, err := nameType.GetSimilarMatch(name, preloadTable)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Couldn't find canonical entity"})
		return
	}

	// Return successful response
	c.JSON(200, canonicalEntity)
	return
}

// UpdateName updates name by id
func UpdateName(c *gin.Context) {
	// Convert id string into int
	param := c.Params.ByName("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Id": "error on converting id"})
		return
	}
	// Get the name by id
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

	// Name is passed by middlewares
	nameValue, ok := c.Get("name")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "name is not present on middlewares"})
		return
	}

	// Parse nameValue into models.NameTypeInput
	var input models.NameType
	input, ok = nameValue.(models.NameType)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "failed to parse name"})
		return
	}

	// Check if input is the same as the name in the database
	if input.Name == name.Name && input.Classification == name.Classification && input.Metaphone == name.Metaphone && input.NameVariations == name.NameVariations {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Every item on json is the same on the database id " + param})
		return
	} else {
		// Update the name properties if they have changed
		if input.Name != "" && input.Name != name.Name {
			name.Name = input.Name
		}
		if input.Classification != "" && input.Classification != name.Classification {
			name.Classification = input.Classification
		}
		if input.Metaphone != "" && input.Metaphone != name.Metaphone {
			name.Metaphone = input.Metaphone
		}
		if input.NameVariations != "" && input.NameVariations != name.NameVariations {
			name.NameVariations = input.NameVariations
		}

		// Save the updated name to the database
		if err := db.Save(name).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err})
			return
		}

		// Return the updated name
		c.JSON(http.StatusOK, name)

		// Clear the cache
		clearCache(c)

		return
	}
}

// DeleteName deletes name off database by id
func DeleteName(c *gin.Context) {
	var name models.NameType
	// Convert id string into int
	param := c.Params.ByName("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error on converting id"})
		return
	}

	// Check if the name with given id exists
	n, _, err := name.GetNameById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}
	if n.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}

	// Delete the name from the database
	if _, err = name.DeleteNameById(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}
}

// checkCache retrieves the cached name types from the context if they exist, otherwise it retrieves them from the database
func checkCache(c *gin.Context) []models.NameType {
	// Initialize a variable for the NameType type
	var nameType models.NameType
	// Initialize a variable for the cached name types
	var preloadTable []models.NameType
	// Get the cached name types from the context
	cache, existKey := c.Get("nameTypes")
	// If the name types are found in the cache, set them to the `preloadTable` variable
	if existKey {
		preloadTable = cache.([]models.NameType)
	} else {
		// If the name types are not found in the cache, retrieve them from the database
		allNames, err := nameType.GetAllNames()
		if err != nil {
			// If there is an error retrieving the name types from the database, return nil
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "Error on caching all name types"})
			return nil
		}
		// Set the retrieved name types in the cache
		preloadTable = allNames
		c.Set("nameTypes", preloadTable)
	}

	// Return the cached name types
	return preloadTable
}

// clearCache deletes the cached name types from the context if they exist
func clearCache(c *gin.Context) {
	// Get the cached name types from the context
	cache, exist := c.Get("nameTypes")
	// If the name types are found in the cache, delete them
	if exist {
		// Convert the cache to a sync.Map so that we can delete the cached name types
		if cm, ok := cache.(sync.Map); ok {
			cm.Delete("preloadTable")
		}
	}
}
