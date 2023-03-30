package main

import (
	"fmt"
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/log"
	"github.com/Darklabel91/API_Names/routes"
)

const FILENAME = "logs.txt"

func main() {
	r := database.InitDb()
	if r == nil {
		return
	}

	fmt.Println("-	Listening and serving")
	go log.StartExportLog(FILENAME)
	routes.HandleRequests()
}
