package models

import (
	"gorm.io/gorm"
	"time"
)

//NameType main struct
type NameType struct {
	gorm.Model
	Name           string `gorm:"unique" json:"Name,omitempty"`
	Classification string `json:"Classification,omitempty"`
	Metaphone      string `json:"Metaphone,omitempty"`
	NameVariations string `json:"NameVariations,omitempty"`
}

//User is the struct of API users
type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

//NameLevenshtein struct for organizing name variations by Levenshtein
type NameLevenshtein struct {
	Name        string
	Levenshtein float32
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

//InputBody struct for validation middleware
type InputBody struct {
	Email    string `json:"Email,omitempty"`
	Password string `json:"Password,omitempty"`
}
