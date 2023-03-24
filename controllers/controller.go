package controllers

import (
	"errors"
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

//GetID read name by id
func GetID(c *gin.Context) {
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
	database.Db.Raw("select * from name_types where name = ?", strings.ToUpper(n)).Find(&name)

	if name.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name not found"})
		return
	}

	c.JSON(http.StatusOK, name)
	return
}

//SearchSimilarNames search for all similar names by metaphone and Levenshtein method
func SearchSimilarNames(c *gin.Context) {
	//Name to be searched
	name := c.Params.ByName("name")

	var names []models.NameType
	database.Db.Raw("select * from name_types").Find(&names)

	var canonicalEntity models.NameType
	database.Db.Raw("select * from name_types where name = ?", strings.ToUpper(name)).Find(&canonicalEntity)

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

	//build canonical
	if canonicalEntity.ID == 0 {
		ce, err := findCanonical(nameV)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"Not found": err.Error(), "metaphone": mtf})
			return
		}
		canonicalEntity = ce
	}

	//return
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
	c.JSON(200, r)
}

/*-------ALL BELLOW USED ONLY ON searchSimilarNames-------*/

//findCanonical search for every similar name on the database returning the first matched name
func findCanonical(similarNames []string) (models.NameType, error) {
	var canonicalEntity models.NameType

	for _, similarName := range similarNames {
		database.Db.Raw("select * from name_types where name = ?", strings.ToUpper(similarName)).Find(&canonicalEntity)
		if canonicalEntity.ID != 0 {
			return canonicalEntity, nil
		}
	}

	return models.NameType{}, errors.New("couldn't find canonical name")
}

//findSimilarNames returns []models.NameVar and if necessary reduces' threshold to a minimum of 0.5
func findSimilarNames(names []models.NameType, name string, threshold float32) ([]models.NameVar, string) {
	similarNames, mtf := findNames(names, name, threshold)

	//in case of empty return the levenshtein constant is downgraded to the minimum of 0.5
	if len(similarNames) == 0 {
		similarNames, _ = findNames(names, name, threshold-0.1)
		if len(similarNames) == 0 {
			similarNames, _ = findNames(names, name, threshold-0.2)
		}
		if len(similarNames) == 0 {
			similarNames, _ = findNames(names, name, threshold-0.3)
		}
	}

	return similarNames, mtf
}

//findNames return []models.NameVar with all similar names and the metaphone code of searched string, called on  findSimilarNames
func findNames(names []models.NameType, name string, threshold float32) ([]models.NameVar, string) {
	var similarNames []models.NameVar

	mtf := metaphone.Pack(name)
	for _, n := range names {
		if metaphone.IsMetaphoneSimilar(mtf, n.Metaphone) {
			similarity := metaphone.SimilarityBetweenWords(strings.ToLower(name), strings.ToLower(n.Name))
			if similarity >= threshold {
				similarNames = append(similarNames, models.NameVar{Name: n.Name, Levenshtein: similarity})
				varWords := strings.Split(n.NameVariations, "|")
				for _, vw := range varWords {
					if vw != "" {
						similarNames = append(similarNames, models.NameVar{Name: vw, Levenshtein: similarity})
					}
				}
			}

		}
	}

	return similarNames, mtf

}

//orderByLevenshtein used to sort an array by Levenshtein and len of the name
func orderByLevenshtein(arr []models.NameVar) []string {
	// creates copy of original array
	sortedArr := make([]models.NameVar, len(arr))
	copy(sortedArr, arr)

	// order by func
	sort.Slice(sortedArr, func(i, j int) bool {
		if sortedArr[i].Levenshtein != sortedArr[j].Levenshtein {
			return sortedArr[i].Levenshtein > sortedArr[j].Levenshtein
		} else {
			return len(sortedArr[i].Name) < len(sortedArr[j].Name)
		}
	})

	//return array
	var retArr []string
	for _, lv := range sortedArr {
		retArr = append(retArr, lv.Name)
	}

	//return without duplicates
	return removeDuplicates(retArr)
}

//removeDuplicates remove duplicates of []string, called on orderByLevenshtein
func removeDuplicates(arr []string) []string {
	var cleanArr []string

	for _, a := range arr {
		if !contains(cleanArr, a) {
			cleanArr = append(cleanArr, a)
		}
	}

	return cleanArr
}

//contains verifies if []string already has a specific string, called on removeDuplicates
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
