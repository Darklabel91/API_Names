package models

import (
	"sort"
)

//NameLevenshtein struct for organizing name variations by Levenshtein
type NameLevenshtein struct {
	Name        string
	Levenshtein float32
}

const LEVENSHTEIN = 0.8

//OrderByLevenshtein used to sort an array by Levenshtein and len of the name
func (*NameLevenshtein) OrderByLevenshtein(arr []NameLevenshtein) []string {
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

//removeDuplicates remove duplicates of []string, called on OrderByLevenshtein
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
