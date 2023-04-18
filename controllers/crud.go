package controllers

import (
	"fmt"
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error on getting name from middlewares"})
		return
	}

	// Parse nameValue into models.NameTypeInput
	newName, ok := nameValue.(models.NameType)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to parse name"})
		return
	}

	// Check the cache
	preloadTable := checkCache(c)

	// Check if there's an exact name on the database
	for _, name := range preloadTable {
		if name.Name == newName.Name {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "name already on the database"})
			return
		}
	}

	fmt.Println("chegou aqui")

	// Create name
	err := newName.CreateName()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error on creating name"})
		return
	}

	// Clear cache
	clearCache(c)

	// Return successful response
	c.JSON(http.StatusOK, gin.H{"Message": "Name created"})

	// Clear the cache
	clearCache(c)

	return
}

//GetID reads a name by id
func GetID(c *gin.Context) {
	// Convert id string into int
	param := c.Params.ByName("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error parsing id"})
		return
	}

	// Get the name by id
	n, _, err := models.GetNameById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error getting name by id"})
		return
	}
	if n.ID == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "name id not found"})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, n)

	return
}

//GetName reads a name by name
func GetName(c *gin.Context) {
	// Get name to be searched
	param := c.Params.ByName("name")

	// Search for name
	n, err := models.GetNameByName(strings.ToUpper(param))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error getting name by name"})
		return
	}
	if n.ID == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "name not found"})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, n)
	return
}

//GetMetaphoneMatch reads a name by metaphone
func GetMetaphoneMatch(c *gin.Context) {
	// Get name to be searched
	name := c.Params.ByName("name")

	// Check the cache
	preloadTable := checkCache(c)

	// Search for similar names
	canonicalEntity, err := models.GetSimilarMatch(name, preloadTable)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "error finding canonical entity"})
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error on parsing id"})
		return
	}

	// Get the name by id
	name, db, err := models.GetNameById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error getting name by id"})
		return
	}

	// Name is passed by middlewares
	nameValue, ok := c.Get("name")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "name is not in middleware context"})
		return
	}

	// Parse nameValue into models.NameTypeInput
	updateName, ok := nameValue.(models.NameType)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to parse name"})
		return
	}

	// Update the name get by id with the updated struct
	un, err := name.UpdateName(db, updateName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error on updating item: " + err.Error()})
		return
	}

	// Return the updated name
	c.JSON(http.StatusOK, un)

	// Clear the cache
	clearCache(c)

	return
}

// DeleteName deletes name off database by id
func DeleteName(c *gin.Context) {
	// Convert id string into int
	param := c.Params.ByName("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error on converting id"})
		return
	}

	// Check if the name with given id exists
	name, _, err := models.GetNameById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "error on deleting name: " + err.Error()})
		return
	}

	// Delete the name from the database
	err = name.DeleteName()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "error on deleting name: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "id deleted"})

	// Clear the cache
	clearCache(c)

	return
}

// checkCache retrieves the cached name types from the context if they exist, otherwise it retrieves them from the database
func checkCache(c *gin.Context) []models.NameType {
	// Initialize a variable for the cached name types
	var preloadTable []models.NameType

	// Get the cached name types from the context
	cache, existKey := c.Get("nameTypes")
	if existKey {
		preloadTable = cache.([]models.NameType)
	} else {
		allNames, err := models.GetAllNames()
		if err != nil {
			// If there is an error retrieving the name types from the database, return nil
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error on caching all name types"})
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
