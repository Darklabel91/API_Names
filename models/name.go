package models

import (
	"errors"
	"github.com/Darklabel91/metaphone-br"
	"gorm.io/gorm"
	"strings"
)

var DB *gorm.DB
var IPs []string

type NameType struct {
	gorm.Model
	Name           string `gorm:"unique" json:"Name,omitempty"`
	Classification string `json:"Classification,omitempty"`
	Metaphone      string `json:"Metaphone,omitempty"`
	NameVariations string `json:"NameVariations,omitempty"`
}

func (n *NameType) CreateName() (*NameType, error) {
	name := n
	r := DB.Create(&name)
	if r.Error != nil {
		return nil, r.Error
	}
	return name, nil
}

func (*NameType) GetAllNames() ([]NameType, error) {
	var Names []NameType
	r := DB.Raw("select * from name_types").Find(&Names)
	if r.Error != nil {
		return nil, r.Error
	}
	return Names, nil
}

func (*NameType) GetNameById(id int) (*NameType, *gorm.DB, error) {
	var getName NameType
	data := DB.Raw("select * from name_types where id = ?", id).Find(&getName)
	if data.Error != nil {
		return nil, nil, data.Error
	}
	return &getName, data, nil
}

func (*NameType) GetNameByName(name string) (*NameType, error) {
	var getName NameType
	data := DB.Raw("select * from name_types where name = ?", name).Find(&getName)
	if data.Error != nil {
		return nil, data.Error
	}
	return &getName, nil
}

func (*NameType) GetNameByMetaphone(mtf string) ([]NameType, error) {
	var getNames []NameType
	data := DB.Raw("select * from name_types where metaphone = ?", mtf).Find(&getNames)
	if data.Error != nil {
		return nil, data.Error
	}
	return getNames, nil
}

func (n *NameType) GetSimilarMatch(name string, allNames []NameType) (*NameType, error) {
	//search perfect match on database
	perfectMatch, err := n.GetNameByName(strings.ToUpper(name))
	if err != nil {
		return nil, err
	}
	if perfectMatch.ID != 0 {
		return perfectMatch, nil
	}

	//Search by similar match
	nameMetaphone := metaphone.Pack(name)

	//1- get the exact metaphone match
	metaphoneNameMatches := n.SearchCacheMetaphone(nameMetaphone, allNames)
	//case we don't find the exact match of metaphone we search all similar metaphone
	if len(metaphoneNameMatches) == 0 {
		//get all similar metaphone code
		metaphoneNameMatches = n.SearchSimilarMetaphone(nameMetaphone, allNames)
		if len(metaphoneNameMatches) == 0 {
			return nil, errors.New("no metaphone matches found")
		}
	}

	//2- get all similar names by metaphone list
	similarNames := n.SearchSimilarNames(name, metaphoneNameMatches, LEVENSHTEIN)
	if len(similarNames) == 0 {
		return nil, errors.New("no similar name matches")
	}
	//case similarNames is too small we search for all similar names of all similar names listed so far
	if len(similarNames) < 5 {
		for _, sn := range similarNames {
			similar := n.SearchSimilarNames(sn.Name, metaphoneNameMatches, LEVENSHTEIN)
			similarNames = append(similarNames, similar...)
		}
	}

	//3- order all similarNames by LEVENSHTEIN from high to low
	var nameLevenshtein NameLevenshtein
	similarNamesOrderedByLevenshtein := nameLevenshtein.OrderByLevenshtein(similarNames)

	//4- return the canonical name combined with similar names ordered by levenshtein
	canonicalEntity, err := n.SearchCanonicalName(name, LEVENSHTEIN, allNames, metaphoneNameMatches, similarNamesOrderedByLevenshtein)
	if err != nil {
		return nil, err
	}

	return canonicalEntity, nil

}

func (*NameType) DeleteNameById(id int) (NameType, error) {
	var getName NameType
	r := DB.Raw("select * from name_types where id = ?", id).Find(&getName)
	if r.Error != nil {
		return NameType{}, r.Error
	}
	return getName, nil
}

func (*NameType) SearchSimilarMetaphone(paradigmMetaphone string, allNames []NameType) []NameType {
	var returnNames []NameType
	for _, name := range allNames {
		if metaphone.IsMetaphoneSimilar(paradigmMetaphone, name.Metaphone) {
			returnNames = append(returnNames, name)
		}
	}

	return returnNames
}

func (*NameType) SearchSimilarNames(paradigmName string, allNames []NameType, threshold float32) []NameLevenshtein {
	var similarNames []NameLevenshtein

	for _, name := range allNames {
		similarity := metaphone.SimilarityBetweenWords(strings.ToLower(paradigmName), strings.ToLower(name.Name))
		if similarity >= threshold {
			similarNames = append(similarNames, NameLevenshtein{Name: name.Name, Levenshtein: similarity})
			varWords := strings.Split(name.NameVariations, "|")
			for _, vw := range varWords {
				if vw != "" {
					similarNames = append(similarNames, NameLevenshtein{Name: vw, Levenshtein: similarity})
				}
			}
		}
	}

	//reduce the threshold in 0.1 if noting is found
	if len(similarNames) == 0 {
		for _, name := range allNames {
			similarity := metaphone.SimilarityBetweenWords(strings.ToLower(paradigmName), strings.ToLower(name.Name))
			if similarity >= threshold-0.1 {
				similarNames = append(similarNames, NameLevenshtein{Name: name.Name, Levenshtein: similarity})
				varWords := strings.Split(name.NameVariations, "|")
				for _, vw := range varWords {
					if vw != "" {
						similarNames = append(similarNames, NameLevenshtein{Name: vw, Levenshtein: similarity})
					}
				}
			}
		}
		return similarNames
	}

	return similarNames
}

func (*NameType) SearchCanonicalName(paradigmName string, threshold float32, allNames []NameType, matchingMetaphoneNames []NameType, nameVariations []string) (*NameType, error) {
	n := strings.ToUpper(paradigmName)

	//transform the nameVariations into a string to be returned
	var rNv string
	for _, nv := range nameVariations {
		rNv += nv + " | "
	}

	//search exact match on matchingMetaphoneNames
	for _, similarName := range matchingMetaphoneNames {
		if similarName.Name == n {
			similarName.NameVariations = rNv
			return &similarName, nil
		}
	}

	//search for similar names on matchingMetaphoneNames
	for _, similarName := range matchingMetaphoneNames {
		sn := strings.ToUpper(similarName.NameVariations)
		if metaphone.SimilarityBetweenWords(n, sn) >= threshold {
			similarName.NameVariations = rNv
			return &similarName, nil
		}
	}

	//search exact match on nameVariations
	for _, similarName := range nameVariations {
		sn := strings.ToUpper(similarName)
		if sn == n {
			for _, name := range allNames {
				if name.Name == n {
					name.NameVariations = rNv
					return &name, nil
				}
			}
		}
	}

	//search for similar names on nameVariations
	for _, similarName := range nameVariations {
		sn := strings.ToUpper(similarName)
		if metaphone.SimilarityBetweenWords(n, sn) >= threshold {
			for _, name := range allNames {
				if name.Name == sn {
					name.NameVariations = rNv
					return &name, nil
				}
			}
		}
	}

	//case none are found we establish a return similarity for names 0.1 bellow the original threshold
	for _, similarName := range nameVariations {
		sn := strings.ToUpper(similarName)
		if metaphone.SimilarityBetweenWords(n, sn) >= threshold-0.1 {
			for _, name := range allNames {
				if name.Name == sn {
					name.NameVariations = rNv
					return &name, nil
				}
			}
		}
	}

	return &NameType{}, errors.New("couldn't find canonical name")
}

func (*NameType) SearchCacheName(name string, cache []NameType) (*NameType, bool) {
	for _, c := range cache {
		if c.Name == name {
			return &c, true
		}
	}

	return &NameType{}, false

}

func (*NameType) SearchCacheMetaphone(metaphone string, cache []NameType) []NameType {
	var nameTypes []NameType
	for _, c := range cache {
		if c.Metaphone == metaphone {
			nameTypes = append(nameTypes, c)
		}
	}

	if len(nameTypes) != 0 {
		return nameTypes
	}

	return nil
}
