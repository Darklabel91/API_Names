package controllers

import (
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/metaphone"
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
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

//SearchSimilarNames search for all similar names by metaphone and levenshtein method
func SearchSimilarNames(c *gin.Context) {
	var names []models.NameType

	//Name to be searched
	name := c.Params.ByName("name")
	database.Db.Find(&names)

	mtf := metaphone.Pack(name)
	var similarNames []models.NameVar
	for _, n := range names {
		if metaphone.IsMetaphoneSimilar(mtf, n.Metaphone) {
			smlt := metaphone.SimilarityBetweenWords(strings.ToLower(name), strings.ToLower(n.Name))
			if smlt >= levenshtein {
				similarNames = append(similarNames, models.NameVar{Name: n.Name, Levenshtein: smlt})
				varWords := strings.Split(n.NameVariations, "|")
				for _, vw := range varWords {
					if vw != "" {
						similarNames = append(similarNames, models.NameVar{Name: vw, Levenshtein: smlt})
					}
				}
			}

			if len(similarNames) == 0 {
				similarNames = append(similarNames, models.NameVar{Name: n.Name, Levenshtein: smlt})
				varWords := strings.Split(n.NameVariations, "|")
				for _, vw := range varWords {
					if vw != "" {
						similarNames = append(similarNames, models.NameVar{Name: vw, Levenshtein: smlt})
					}
				}
			}
		}
	}

	if len(names) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "metaphone not found", "metaphone": mtf})
		return
	}

	nameV := orderByLevenshtein(similarNames)

	c.JSON(200, gin.H{
		"Name":           strings.ToUpper(name),
		"metaphone":      mtf,
		"nameVariations": nameV,
	})

}

//orderByLevenshtein used to sort an array by Levenshtein
func orderByLevenshtein(arr []models.NameVar) []string {
	// creates copy of original array
	sortedArr := make([]models.NameVar, len(arr))
	copy(sortedArr, arr)

	// compilation func
	cmp := func(i, j int) bool {
		return sortedArr[i].Levenshtein > sortedArr[j].Levenshtein
	}

	// order by func
	sort.Slice(sortedArr, cmp)

	var retArry []string
	for _, lv := range sortedArr {
		if lv.Levenshtein != float32(0) {
			retArry = append(retArry, lv.Name)
		}

	}

	return retArry
}
