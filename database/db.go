package database

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/Darklabel91/API_Names/models"
	"github.com/Darklabel91/metaphone-br"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"os"
	"time"
)

var DB *gorm.DB

//ConnectDB open connection and migrate tables ORM
func ConnectDB() *gorm.DB {
	//load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file was not found. You should add a .env file on project root with:\nDB_USERNAME \nDB_PASSWORD \nDB_NAME \nDB_HOST \nDB_PORT \nSECRET")
		return nil
	}

	//get .env variables
	var (
		DbUsername = os.Getenv("DB_USERNAME")
		DbPassword = os.Getenv("DB_PASSWORD")
		DbName     = os.Getenv("DB_NAME")
		DbHost     = os.Getenv("DB_HOST")
		DbPort     = os.Getenv("DB_PORT")
	)

	//create database
	err = createDatabase(DbHost, DbUsername, DbPassword, DbName)
	if err != nil {
		fmt.Printf("Error on gorm creating database : error=%v\n", err)
		return nil
	}

	dsn := DbUsername + ":" + DbPassword + "@tcp" + "(" + DbHost + ":" + DbPort + ")/" + DbName + "?" + "parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error connecting to database : error=%v\n", err)
		return nil
	}

	err = db.AutoMigrate(models.NameType{}, models.User{}, models.Log{})
	if err != nil {
		fmt.Printf("Error on gorm auto migrate to database : error=%v\n", err)
		return nil
	}

	DB = db

	err = uploadCSVNameTypes()
	if err != nil {
		fmt.Printf("Error uploading .csv to database : error=%v\n", err)
		return nil
	}

	return db
}

//createDatabase runs create database script
func createDatabase(host, username, password, dbName string) error {
	// Set up the MySQL DSN string
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", username, password, host)

	// Open a connection to the MySQL server
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	// Check if the database already exists
	var result int64
	db.Raw("SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?", dbName).Scan(&result)
	if result > 0 {
		return nil
	}

	// Create the database
	err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}
	fmt.Println("-	Create Database")

	return nil
}

//UploadCSVNameTypes upload the .csv file on database folder on names table
func uploadCSVNameTypes() error {
	var name models.NameType
	DB.Raw("select * from name_types where id = 1").Find(&name)

	if name.ID == 0 {
		start := time.Now()
		fmt.Println("-	Upload data start")

		filePath := "database/name_types .csv"
		file, err := os.Open(filePath)
		if err != nil {
			return errors.New("Error opening file:" + err.Error())

		}
		defer file.Close()

		reader := csv.NewReader(file)
		var rows [][]string
		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return errors.New("error reading CSV:" + err.Error())
			}
			rows = append(rows, row)
		}

		for i, row := range rows {
			if i != 0 {
				nameType := models.NameType{
					Name:           row[0],
					Classification: row[1],
					Metaphone:      metaphone.Pack(row[0]),
					NameVariations: row[3],
				}
				if err = DB.Create(&nameType).Error; err != nil {
					return errors.New("error creating NameType:" + err.Error())
				}
			}
		}

		fmt.Println("-	Upload data finished", time.Since(start).String())
		return nil
	}
	return nil

}
