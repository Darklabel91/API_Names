package models

import (
	"errors"
	"sort"
)

// NameSimilarity contains a name and a Levenshtein score
type NameSimilarity struct {
	Name       string
	Similarity float32
}

const SimilarityThreshold = 0.8

// OrderBySimilarity sorts an array of NameSimilarity objects by descending similarity
// and then by ascending name length.
func OrderBySimilarity(arr []NameSimilarity) ([]string, error) {
	if len(arr) == 0 {
		return nil, errors.New("input array is empty")
	}

	// create a copy of the original array
	sortedArr := make([]NameSimilarity, len(arr))
	copy(sortedArr, arr)

	// sort the array by similarity and then by name length
	sort.Slice(sortedArr, func(i, j int) bool {
		if sortedArr[i].Similarity != sortedArr[j].Similarity {
			return sortedArr[i].Similarity > sortedArr[j].Similarity
		} else {
			return len(sortedArr[i].Name) < len(sortedArr[j].Name)
		}
	})

	// convert the sorted array to an array of strings
	var retArr []string
	for _, sim := range sortedArr {
		retArr = append(retArr, sim.Name)
	}

	// remove duplicates from the sorted array
	retArr = removeDuplicates(retArr)

	return retArr, nil
}

// removeDuplicates removes duplicate strings from a slice of strings.
func removeDuplicates(arr []string) []string {
	var cleanArr []string

	for _, a := range arr {
		if !contains(cleanArr, a) {
			cleanArr = append(cleanArr, a)
		}
	}

	return cleanArr
}

// contains returns true if a slice of strings contains a given string, false otherwise.
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
