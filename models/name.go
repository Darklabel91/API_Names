package models

import (
	"errors"
	"fmt"
	"github.com/Darklabel91/metaphone-br"
	"gorm.io/gorm"
	"strings"
)

var DB *gorm.DB
var IPs []string

// NameType is a struct representing a name record
type NameType struct {
	gorm.Model
	Name           string `gorm:"unique" json:"Name,omitempty"`
	Classification string `json:"Classification,omitempty"`
	Metaphone      string `json:"Metaphone,omitempty"`
	NameVariations string `json:"NameVariations,omitempty"`
}

// CreateName creates a new name record
func (n *NameType) CreateName() (*NameType, error) {
	name := n
	r := DB.Create(&name)
	if r.Error != nil {
		return nil, r.Error
	}
	return name, nil
}

// GetAllNames returns all non-deleted names in the database
func (*NameType) GetAllNames() ([]NameType, error) {
	var Names []NameType
	r := DB.Raw("SELECT * FROM name_types WHERE name_types.deleted_at IS NULL").Find(&Names)
	if r.Error != nil {
		return nil, r.Error
	}
	return Names, nil
}

// GetNameById returns the name record with the given ID (non-deleted)
func (*NameType) GetNameById(id int) (*NameType, *gorm.DB, error) {
	var getName NameType
	data := DB.Raw("SELECT * FROM name_types WHERE id = ? AND name_types.deleted_at IS NULL", id).Find(&getName)
	if data.Error != nil {
		return nil, nil, data.Error
	}
	return &getName, data, nil
}

// GetNameByName returns the name record with the given name (non-deleted)
func (*NameType) GetNameByName(name string) (*NameType, error) {
	var getName NameType
	data := DB.Raw("SELECT * FROM name_types WHERE name = ? AND name_types.deleted_at IS NULL", name).Find(&getName)
	if data.Error != nil {
		return nil, data.Error
	}
	return &getName, nil
}

// GetNameByMetaphone returns all non-deleted name records with the given metaphone
func (*NameType) GetNameByMetaphone(mtf string) ([]NameType, error) {
	var getNames []NameType
	data := DB.Raw("SELECT * FROM name_types WHERE metaphone = ? AND name_types.deleted_at IS NULL", mtf).Find(&getNames)
	if data.Error != nil {
		return nil, data.Error
	}
	return getNames, nil
}

// GetSimilarMatch searches for a similar match for a given name in a slice of NameType.
func (n *NameType) GetSimilarMatch(name string, allNames []NameType) (*NameType, error) {
	// Search for an exact match in the database.
	perfectMatch, err := n.GetNameByName(strings.ToUpper(name))
	if err != nil {
		return nil, fmt.Errorf("failed to search for exact match: %w", err)
	}
	if perfectMatch.ID != 0 {
		return perfectMatch, nil
	}

	// Search for a similar match.
	nameMetaphone := metaphone.Pack(name)

	// Search for the exact metaphone match.
	exactMetaphoneMatches := n.SearchCacheMetaphone(nameMetaphone, allNames)
	if len(exactMetaphoneMatches) == 0 {
		// Search for all similar metaphone codes if no exact match is found.
		exactMetaphoneMatches = n.SearchSimilarMetaphone(nameMetaphone, allNames)
		if len(exactMetaphoneMatches) == 0 {
			return nil, fmt.Errorf("no matches found for name %q", name)
		}
	}

	// Get all similar names by metaphone list.
	similarNames := n.SearchSimilarNames(name, exactMetaphoneMatches, SimilarityThreshold)
	if len(similarNames) == 0 {
		return nil, fmt.Errorf("no similar names found for %q", name)
	}

	// Search for all similar names of all similar names listed so far if similarNames is too small.
	if len(similarNames) < 5 {
		for _, sn := range similarNames {
			similar := n.SearchSimilarNames(sn.Name, exactMetaphoneMatches, SimilarityThreshold)
			similarNames = append(similarNames, similar...)
		}
	}

	// Order all similarNames by LEVENSHTEIN from high to low.
	similarNamesOrderedByLevenshtein, err := OrderBySimilarity(similarNames)
	if err != nil {
		return nil, fmt.Errorf("failed to order by similar names: %w", err)
	}

	// Return the canonical name combined with similar names ordered by levenshtein.
	canonicalEntity, err := n.SearchCanonicalName(name, SimilarityThreshold, allNames, exactMetaphoneMatches, similarNamesOrderedByLevenshtein)
	if err != nil {
		return nil, fmt.Errorf("failed to search for canonical name: %w", err)
	}

	return canonicalEntity, nil
}

// DeleteNameById deletes a name from the database by its ID.
func (*NameType) DeleteNameById(id int) (NameType, error) {
	var getName NameType
	r := DB.Where("id = ?", id).Delete(&getName)
	if r.Error != nil {
		return NameType{}, fmt.Errorf("failed to delete name with ID %d: %w", id, r.Error)
	}
	return getName, nil
}

// SearchSimilarMetaphone returns a slice of NameType elements that have a metaphone similar to the given paradigmMetaphone
func (*NameType) SearchSimilarMetaphone(paradigmMetaphone string, allNames []NameType) []NameType {
	// create an empty slice to store the return values
	var returnNames []NameType
	// iterate over all the names in allNames
	for _, name := range allNames {
		// if the metaphone of the name is similar to the given paradigmMetaphone
		if metaphone.IsMetaphoneSimilar(paradigmMetaphone, name.Metaphone) {
			// add the name to the returnNames slice
			returnNames = append(returnNames, name)
		}
	}
	// return the slice of NameType elements that have a metaphone similar to the given paradigmMetaphone
	return returnNames
}

// SearchSimilarNames returns a slice of NameLevenshtein elements that have a similarity score higher than the given threshold to the given paradigmName
func (*NameType) SearchSimilarNames(paradigmName string, allNames []NameType, threshold float32) []NameSimilarity {
	// create an empty slice to store the return values
	var similarNames []NameSimilarity
	// iterate over all the names in allNames
	for _, name := range allNames {
		// calculate the similarity between the paradigmName and the name using the metaphone package
		similarity := metaphone.SimilarityBetweenWords(strings.ToLower(paradigmName), strings.ToLower(name.Name))
		// if the similarity score is higher than the given threshold
		if similarity >= threshold {
			// create a new NameLevenshtein element with the name and the similarity score
			similarName := NameSimilarity{Name: name.Name, Similarity: similarity}
			// add the new NameLevenshtein element to the similarNames slice
			similarNames = append(similarNames, similarName)
			// split the name variations string of the name and iterate over each variation
			for _, vw := range strings.Split(name.NameVariations, "|") {
				// if the variation is not empty
				if vw != "" {
					// create a new NameLevenshtein element with the variation and the similarity score
					variationName := NameSimilarity{Name: vw, Similarity: similarity}
					// add the new NameLevenshtein element to the similarNames slice
					similarNames = append(similarNames, variationName)
				}
			}
		}
	}
	// if no names were found with a similarity score higher than the threshold, try again with a lower threshold
	if len(similarNames) == 0 {
		for _, name := range allNames {
			similarity := metaphone.SimilarityBetweenWords(strings.ToLower(paradigmName), strings.ToLower(name.Name))
			if similarity >= threshold-0.1 {
				similarName := NameSimilarity{Name: name.Name, Similarity: similarity}
				similarNames = append(similarNames, similarName)
				for _, vw := range strings.Split(name.NameVariations, "|") {
					if vw != "" {
						variationName := NameSimilarity{Name: vw, Similarity: similarity}
						similarNames = append(similarNames, variationName)
					}
				}
			}
		}
		return similarNames
	}
	return similarNames
}

// SearchCanonicalName searches for a canonical name in a list of names using a given threshold for similarity matching.
func (*NameType) SearchCanonicalName(paradigmName string, threshold float32, allNames []NameType, matchingMetaphoneNames []NameType, nameVariations []string) (*NameType, error) {
	// Convert the input name to uppercase.
	n := strings.ToUpper(paradigmName)

	// Transform the nameVariations into a string to be returned.
	var rNv string
	for _, nv := range nameVariations {
		rNv += nv + " | "
	}

	// Search for exact match on matchingMetaphoneNames.
	for _, similarName := range matchingMetaphoneNames {
		if similarName.Name == n {
			// If an exact match is found, add the name variations and return the result.
			similarName.NameVariations = rNv
			return &similarName, nil
		}
	}

	// Search for similar names on matchingMetaphoneNames.
	for _, similarName := range matchingMetaphoneNames {
		// Convert the name variations to uppercase for comparison.
		sn := strings.ToUpper(similarName.NameVariations)
		if metaphone.SimilarityBetweenWords(n, sn) >= threshold {
			// If a similar name is found, add the name variations and return the result.
			similarName.NameVariations = rNv
			return &similarName, nil
		}
	}

	// Search for exact match on nameVariations.
	for _, similarName := range nameVariations {
		// Convert the name variations to uppercase for comparison.
		sn := strings.ToUpper(similarName)
		if sn == n {
			// If an exact match is found, search for the corresponding name in allNames, add the name variations, and return the result.
			for _, name := range allNames {
				if name.Name == n {
					name.NameVariations = rNv
					return &name, nil
				}
			}
		}
	}

	// Search for similar names on nameVariations.
	for _, similarName := range nameVariations {
		// Convert the name variations to uppercase for comparison.
		sn := strings.ToUpper(similarName)
		if metaphone.SimilarityBetweenWords(n, sn) >= threshold {
			// If a similar name is found, search for the corresponding name in allNames, add the name variations, and return the result.
			for _, name := range allNames {
				if name.Name == sn {
					name.NameVariations = rNv
					return &name, nil
				}
			}
		}
	}

	// If none of the above searches succeed, search for similar names on nameVariations with a lower threshold.
	for _, similarName := range nameVariations {
		// Convert the name variations to uppercase for comparison.
		sn := strings.ToUpper(similarName)
		if metaphone.SimilarityBetweenWords(n, sn) >= threshold-0.1 {
			// If a similar name is found, search for the corresponding name in allNames, add the name variations, and return the result.
			for _, name := range allNames {
				if name.Name == sn {
					name.NameVariations = rNv
					return &name, nil
				}
			}
		}
	}

	// If no match is found, return an error.
	return &NameType{}, errors.New("couldn't find canonical name")
}

// SearchCacheName searches for a name in the cache and returns a pointer to its corresponding NameType object
// along with a boolean indicating whether the name was found or not
func (*NameType) SearchCacheName(name string, cache []NameType) (*NameType, bool) {
	// iterate over the cache and look for a name that matches the given name
	for _, c := range cache {
		if c.Name == name {
			// if the name is found, return a pointer to the NameType object and true
			return &c, true
		}
	}

	// if the name is not found, return a pointer to an empty NameType object and false
	return &NameType{}, false
}

// SearchCacheMetaphone searches for all NameType objects in the cache that have a matching metaphone value
// and returns them as a slice
func (*NameType) SearchCacheMetaphone(metaphone string, cache []NameType) []NameType {
	// create an empty slice to hold the NameType objects with matching metaphone values
	var nameTypes []NameType
	// iterate over the cache and look for NameType objects with metaphone values that match the given metaphone value
	for _, c := range cache {
		if c.Metaphone == metaphone {
			// if a match is found, append the NameType object to the slice
			nameTypes = append(nameTypes, c)
		}
	}

	// if at least one NameType object is found, return the slice
	if len(nameTypes) != 0 {
		return nameTypes
	}

	// if no NameType object is found, return nil
	return nil
}
