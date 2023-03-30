package models

import (
	"gorm.io/gorm"
	"time"
)

//NameLevenshtein struct for organizing name variations by Levenshtein
type NameLevenshtein struct {
	Name        string
	Levenshtein float32
}

//InputBody struct for validation middleware
type InputBody struct {
	Email    string `json:"Email,omitempty"`
	Password string `json:"Password,omitempty"`
}

//MetaphoneR only use for SearchSimilarNames return
type MetaphoneR struct {
	ID             uint           `json:"ID,omitempty"`
	CreatedAt      time.Time      `json:"CreatedAt,omitempty"`
	UpdatedAt      time.Time      `json:"UpdatedAt,omitempty"`
	DeletedAt      gorm.DeletedAt `json:"DeletedAt,omitempty"`
	Name           string         `json:"Name,omitempty"`
	Classification string         `json:"Classification,omitempty"`
	Metaphone      string         `json:"Metaphone,omitempty"`
	NameVariations []string       `json:"NameVariations,omitempty"`
}
