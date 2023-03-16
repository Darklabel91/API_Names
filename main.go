package main

import (
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/routes"
)

func main() {
	database.InitDb()
	routes.HandleRequests()
}
