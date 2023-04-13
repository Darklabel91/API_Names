package main

import (
	"log"

	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/models"
	"github.com/Darklabel91/API_Names/routes"
)

func init() {
	// Connect to the database.
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", db.Error)
	}
	models.DB = db

	// Create the root user.
	if err := models.CreateRoot(); err != nil {
		log.Fatalf("Error creating root user: %v", err)
	}

	// Get the list of trusted IPs from the database.
	trustedIPs, err := models.TrustedIPs()
	if err != nil {
		log.Fatalf("Error getting trusted IPs: %v", err)
	}

	// Store the list of trusted IPs in the models package.
	models.IPs = trustedIPs
}

func main() {
	// Handle incoming HTTP requests.
	log.Println("-	Listening and serving...")
	err := routes.HandleRequests()
	if err != nil {
		log.Fatalf("Error handling requests: %v", err)
	}
}
