package controllers

import (
	"errors"
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/models"
	"github.com/Darklabel91/metaphone-br"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strings"
)

const levenshtein = 0.8

//NameLevenshtein struct for organizing name variations by Levenshtein
type NameLevenshtein struct {
	Name        string
	Levenshtein float32
}

//GetSimilarNames search for all similar names by metaphone and Levenshtein method
func GetSimilarNames(c *gin.Context) {
	var metaphoneNames []models.NameType

	//name to be searched
	name := c.Params.ByName("name")
	nameMetaphone := metaphone.Pack(name)

	//Check the cache
	var preloadTable []models.NameType
	cache, existKey := c.Get("nameTypes")
	if existKey {
		preloadTable = cache.([]models.NameType)
	} else {
		if err := database.DB.Find(&preloadTable).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to preload nameTypes"})
			return
		}
	}

	//search perfect match
	nameCache, existName := searchCacheName(name, preloadTable)
	if existName {
		r := models.NameType{
			Name:           nameCache.Name,
			Classification: nameCache.Classification,
			Metaphone:      nameCache.Metaphone,
			NameVariations: nameCache.NameVariations,
		}
		c.JSON(200, r)
		return
	} else {
		//search perfect match on database
		_, err := nameCache.GetNameByName(strings.ToUpper(name))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "error on getting name", "Error": err})
			return
		}
		if len(metaphoneNames) == 1 {
			r := models.NameType{
				Name:           metaphoneNames[0].Name,
				Classification: metaphoneNames[0].Classification,
				Metaphone:      metaphoneNames[0].Metaphone,
				NameVariations: metaphoneNames[0].NameVariations,
			}
			c.JSON(200, r)
			return
		}
	}

	//search metaphone
	var nameType models.NameType
	metaphoneNames, existMetaphone := searchCacheMetaphone(nameMetaphone, preloadTable)
	if !existMetaphone {
		mtn, err := nameType.GetNameByMetaphone(nameMetaphone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "error on get names by metaphone", "Error": err})
			return
		}
		metaphoneNames = mtn
	}

	//find all metaphoneNames matching metaphone
	similarNames := findNames(metaphoneNames, name, levenshtein)

	//for recall purposes we can't only search for metaphone exact match's if no similar word is found
	if len(metaphoneNames) == 0 || len(similarNames) == 0 {
		metaphoneNames = searchForAllSimilarMetaphone(nameMetaphone, preloadTable)
		similarNames = findNames(metaphoneNames, name, levenshtein)

		if len(metaphoneNames) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"Not found": "metaphone not found", "metaphone": nameMetaphone})
			return
		}

		if len(similarNames) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"Not found": "similar names not found", "metaphone": nameMetaphone})
			return
		}
	}

	//when the similar metaphoneNames result's in less than 5 we search for every similar name of all similar metaphoneNames founded previously
	//this step can be ignored if you want to
	if len(similarNames) < 5 {
		for _, n := range similarNames {
			similar := findNames(metaphoneNames, n.Name, levenshtein)
			similarNames = append(similarNames, similar...)
		}
	}

	//order all similar metaphoneNames from high to low Levenshtein
	nameV := orderByLevenshtein(similarNames)

	//finds a name to consider Canonical on the database
	canonicalEntity, err := findCanonical(name, metaphoneNames, nameV)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not found": err.Error(), "metaphone": nameMetaphone})
		return
	}

	var nv string
	for _, variation := range nameV {
		nv += variation + " | "
	}
	r := models.NameType{
		Name:           canonicalEntity.Name,
		Classification: canonicalEntity.Classification,
		Metaphone:      canonicalEntity.Metaphone,
		NameVariations: nv,
	}

	c.JSON(200, r)
	return
}

//searchCacheMetaphone seeks for a given name on cache struct, return the name and a bool true if it is found
func searchCacheName(name string, cache []models.NameType) (models.NameType, bool) {
	for _, c := range cache {
		if c.Name == name {
			return c, true
		}
	}

	return models.NameType{}, false

}

//searchCacheMetaphone seeks for a given name on cache struct, return the name and a bool true if it is found
func searchCacheMetaphone(metaphone string, cache []models.NameType) ([]models.NameType, bool) {
	var nameTypes []models.NameType
	for _, c := range cache {
		if c.Metaphone == metaphone {
			nameTypes = append(nameTypes, c)
		}
	}

	if len(nameTypes) != 0 {
		return nameTypes, true
	}

	return nil, false
}

//----- All needed functions -----//

//searchForAllSimilarMetaphone used in case of not finding exact metaphone match
func searchForAllSimilarMetaphone(mtf string, names []models.NameType) []models.NameType {
	var rNames []models.NameType
	for _, n := range names {
		if metaphone.IsMetaphoneSimilar(mtf, n.Metaphone) {
			rNames = append(rNames, n)
		}
	}

	return rNames
}

//findCanonical search for every similar name on the database returning the first matched name
func findCanonical(name string, matchingMetaphoneNames []models.NameType, nameVariations []string) (models.NameType, error) {
	var canonicalEntity models.NameType
	n := strings.ToUpper(name)

	//search exact match on matchingMetaphoneNames
	for _, similarName := range matchingMetaphoneNames {
		if similarName.Name == n {
			return similarName, nil
		}
	}

	//search for similar names on matchingMetaphoneNames
	for _, similarName := range matchingMetaphoneNames {
		if metaphone.SimilarityBetweenWords(name, similarName.Name) >= levenshtein {
			return similarName, nil
		}
	}

	//search exact match on nameVariations
	for _, similarName := range nameVariations {
		sn := strings.ToUpper(similarName)
		if sn == n {
			ce, err := canonicalEntity.GetNameByName(sn)
			if err != nil {
				return models.NameType{}, err
			}

			if ce.ID != 0 {
				return canonicalEntity, nil
			}
		}
	}

	//in case of failure on other attempts, we search every nameVariations directly on database
	for _, similarName := range nameVariations {
		ce, err := canonicalEntity.GetNameByName(strings.ToUpper(similarName))
		if err != nil {
			return models.NameType{}, err
		}

		if ce.ID != 0 {
			return models.NameType{Name: ce.Name, Classification: ce.Classification, Metaphone: ce.Metaphone, NameVariations: ce.NameVariations}, nil
		}
	}

	return models.NameType{}, errors.New("couldn't find canonical name")
}

//findNames return []NameLevenshtein with all similar names of searched string. For recall purpose we reduce the threshold given in 0.1 in case of empty return
func findNames(names []models.NameType, name string, threshold float32) []NameLevenshtein {
	similarNames := findSimilarNames(name, names, threshold)
	//reduce the threshold given in 0.1 and search again
	if len(similarNames) == 0 {
		similarNames = findSimilarNames(name, names, threshold-0.1)
	}

	return similarNames
}

//findSimilarNames loop for all names given checking the similarity between words by a given threshold, called on findNames
func findSimilarNames(name string, names []models.NameType, threshold float32) []NameLevenshtein {
	var similarNames []NameLevenshtein

	for _, n := range names {
		similarity := metaphone.SimilarityBetweenWords(strings.ToLower(name), strings.ToLower(n.Name))
		if similarity >= threshold {
			similarNames = append(similarNames, NameLevenshtein{Name: n.Name, Levenshtein: similarity})
			varWords := strings.Split(n.NameVariations, "|")
			for _, vw := range varWords {
				if vw != "" {
					similarNames = append(similarNames, NameLevenshtein{Name: vw, Levenshtein: similarity})
				}
			}
		}
	}

	return similarNames
}

//orderByLevenshtein used to sort an array by Levenshtein and len of the name
func orderByLevenshtein(arr []NameLevenshtein) []string {
	// creates copy of original array
	sortedArr := make([]NameLevenshtein, len(arr))
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
