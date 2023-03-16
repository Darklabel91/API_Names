package models

import "gorm.io/gorm"

//var NameCollection []NameType

type NameType struct {
	gorm.Model
	Name           string   `json:"Name,omitempty"`
	Classification string   `json:"Classification,omitempty"`
	Metaphone      string   `json:"Metaphone,omitempty"`
	NameVariations []string `gorm:"type:longText" json:"NameVariations,omitempty"`
}
