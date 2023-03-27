package database

import (
	"fmt"
	"github.com/Darklabel91/API_Names/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var Db *gorm.DB

func InitDb() *gorm.DB {
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	//load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	//get .env variables
	var (
		DbUsername = os.Getenv("DB_USERNAME")
		DbPassword = os.Getenv("DB_PASSWORD")
		DbName     = os.Getenv("DB_NAME")
		DbHost     = os.Getenv("DB_HOST")
		DbPort     = os.Getenv("DB_PORT")
	)

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
