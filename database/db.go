package database

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/Darklabel91/API_Names/models"
	"github.com/Darklabel91/metaphone-br"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"os"
	"time"
)

var Db *gorm.DB

func InitDb() *gorm.DB {
	Db = connectDB()

	err := createRoot()
	if err != nil {
		return nil
	}

	err = uploadCSVNameTypes()
	if err != nil {
		return nil
	}

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

	//create database
	err = createDatabase(DbHost, DbUsername, DbPassword, DbName)
	if err != nil {
		fmt.Printf("Error on gorm creating database : error=%v\n", err)
		return nil
	}

	dsn := DbUsername + ":" + DbPassword + "@tcp" + "(" + DbHost + ":" + DbPort + ")/" + DbName + "?" + "parseTime=true&loc=Local"
	fmt.Println("dsn : ", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error connecting to database : error=%v\n", err)
		return nil
	}

	//create table users
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Printf("Error on gorm auto migrate to database : error=%v\n", err)
		return nil
	}

	//create table name_type
	err = db.AutoMigrate(&models.NameType{})
	if err != nil {
		fmt.Printf("Error on gorm auto migrate to database : error=%v\n", err)
		return nil
	}

	return db
}

func createDatabase(host, username, password, dbName string) error {
	fmt.Println("cra")
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
		fmt.Println("database already exists")
		return nil
	}

	// Create the database
	fmt.Println("creating database")
	err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}

	return nil
}

func createRoot() error {
	var user models.User
	Db.First(&user, 1)

	if user.ID == 0 {
		fmt.Println("creating user root")

		//password
		hash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("SECRET")), 10)
		if err != nil {
			return err
		}

		userRoot := models.User{
			Email:    "root@root.com",
			Password: string(hash),
		}

		Db.Create(&userRoot)

		fmt.Println("root user created")

		return nil
	} else {
		fmt.Println("root user already on the database")
		return nil
	}
}

func uploadCSVNameTypes() error {
	var name models.NameType
	Db.Raw("SELECT * FROM name_types WHERE id = 1").Find(&name)

	if name.ID == 0 {
		start := time.Now()
		fmt.Println("initiating csv upload to the database")

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
				if err = Db.Create(&nameType).Error; err != nil {
					return errors.New("error creating NameType:" + err.Error())
				}
			}
		}

		fmt.Println("csv upload finished in:" + time.Since(start).String())
		return nil
	}
	fmt.Println("csv already imported")
	return nil

}
