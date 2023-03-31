package main

import (
	"fmt"
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/routes"
)

func main() {
	r := database.InitDB()
	if r == nil {
		return
	}

	fmt.Println("-	Listening and serving")
	routes.HandleRequests()
}
