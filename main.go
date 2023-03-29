package main

import (
	"fmt"
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/routes"
)

func main() {
	r := database.InitDb()
	if r == nil {
		return
	}

	fmt.Println("-	live")
	routes.HandleRequests()
}
