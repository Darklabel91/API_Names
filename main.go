package main

import (
	"fmt"
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/routes"
)

func init() {
	database.InitDb()
}

func main() {
	fmt.Println("-	live")
	routes.HandleRequests()
}
