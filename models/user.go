package models

import (
	"gorm.io/gorm"
)

//User is the struct of API users
type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
	IP       string
}
