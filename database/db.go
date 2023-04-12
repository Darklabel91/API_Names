package database

import (
	"encoding/csv"
	"fmt"
	"github.com/Darklabel91/API_Names/models"
	"github.com/Darklabel91/metaphone-br"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"time"
)

// ConnectDB opens a connection to the database and migrates tables using ORM
func ConnectDB() *gorm.DB {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil
	}

	// Get environment variables
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Create database
	err = createDatabase(dbHost, dbUsername, dbPassword, dbName)
	if err != nil {
		return nil
	}

	// Connect to the database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil
	}

	// Migrate tables
	err = db.AutoMigrate(&models.NameType{}, &models.User{}, &models.Log{})
	if err != nil {
		return nil
	}

	// Upload CSV data to NameType table
	err = uploadCSVNameTypes(db)
	if err != nil {
		return nil
	}

	return db
}

// createDatabase runs the create database script
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
	log.Println("-	Created Database")

	return nil
}

// uploadCSVNameTypes the specified CSV file to the database as NameType objects.
func uploadCSVNameTypes(db *gorm.DB) error {
	var name models.NameType
	db.Raw("select * from name_types where id = 1").Find(&name)

	if name.ID == 0 {
		start := time.Now()
		log.Println("-	Upload data start")

		filePath := "database/name_types .csv"
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error opening file:: %v", err)
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
				return fmt.Errorf("error reading file:: %v", err)
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
				if err = db.Create(&nameType).Error; err != nil {
					return fmt.Errorf("error creating NameType:: %v", err)
				}
			}
		}

		log.Println("-	Upload data finished", time.Since(start).String())
		return nil
	}
	return nil
}
