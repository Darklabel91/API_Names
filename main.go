package main

import (
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/routes"
)

func main() {
	//set database
	db := database.InitDB()
	if db == nil {
		return
	}

	//handle requests
	routes.HandleRequests()
}
