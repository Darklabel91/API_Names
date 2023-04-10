package models

//NameLevenshtein struct for organizing name variations by Levenshtein
type NameLevenshtein struct {
	Name        string
	Levenshtein float32
}
