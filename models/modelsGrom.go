package models

import (
	"gorm.io/gorm"
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

//Log is the struct used to register the log file
// Define a struct to hold the log data
type Log struct {
	gorm.Model
	Time    string
	Status  string
	Latency string
	IP      string
	Method  string
	Path    string
}
