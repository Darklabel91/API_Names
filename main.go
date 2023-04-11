package main

import (
	"fmt"
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/models"
	"github.com/Darklabel91/API_Names/routes"
)

func init() {
	db := database.ConnectDB()
	models.DB = db

	err := db.AutoMigrate(models.NameType{}, models.User{}, models.Log{})
	if err != nil {
		fmt.Printf("Error on gorm auto migrate to database : error=%v\n", err)
		return
	}

	err = models.UploadCSVNameTypes()
	if err != nil {
		fmt.Printf("Error on uploading .csv to database : error=%v\n", err)
		return
	}

	err = models.CreateRoot()
	if err != nil {
		fmt.Printf("Error on creating  root user on database: error=%v\n", err)
		return
	}

	IPs, err := models.TrustedIPs()
	if err != nil {
		fmt.Printf("Error on getting the trusted IPs: error=%v\n", err)
		return
	}

	models.IPs = IPs
}

func main() {
	//handle requests
	routes.HandleRequests()
}
