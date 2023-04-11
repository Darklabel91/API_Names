package main

import (
	"fmt"
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/models"
	"github.com/Darklabel91/API_Names/routes"
)

func init() {
	db := database.ConnectDB()
	if db.Error != nil {
		fmt.Printf("Error on connecting to the database : error=%v\n", db.Error)
		return
	}
	models.DB = db

	err := models.CreateRoot()
	if err != nil {
		fmt.Printf("Error on creating  root user on database: error=%v\n", err)
		return
	}

	models.IPs, err = models.TrustedIPs()
	if err != nil {
		fmt.Printf("Error on getting the trusted IPs: error=%v\n", err)
		return
	}
}

func main() {
	//handle requests
	routes.HandleRequests()
}
