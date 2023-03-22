package controllers

import (
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/metaphone"
	"github.com/Darklabel91/API_Names/models"
	"sort"
	"strings"
)

//findCanonical search for every similar name on the database returning the first matched name
func findCanonical(searchName string, similarNames []string) models.NameType {
	var name models.NameType

	database.Db.Where("name = ?", strings.ToUpper(searchName)).Find(&name)
	if name.ID != 0 {
		return name
	}

	for _, similarName := range similarNames {
		database.Db.Where("name = ?", strings.ToUpper(similarName)).Find(&name)
		if name.ID != 0 {
			return name
		}
	}

	return name
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
