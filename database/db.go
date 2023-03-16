package database

import (
	"fmt"
	"github.com/Darklabel91/API_Names/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DbUsername = "root"
	DbPassword = "root"
	DbName     = "namesDatabase"
	DbHost     = "127.0.0.1"
	DbPort     = "3306"
)

var Db *gorm.DB

func InitDb() *gorm.DB {
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	var err error
	dsn := DbUsername + ":" + DbPassword + "@tcp" + "(" + DbHost + ":" + DbPort + ")/" + DbName + "?" + "parseTime=true&loc=Local"
	fmt.Println("dsn : ", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("Error connecting to database : error=%v\n", err)
		return nil
	}

	err = db.AutoMigrate(&models.NameType{})
	if err != nil {
		fmt.Printf("Error on gorm auto migrate to database : error=%v\n", err)
		return nil
	}

	return db
}
