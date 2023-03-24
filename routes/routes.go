package routes

import (
	"github.com/Darklabel91/API_Names/controllers"
	"github.com/gin-gonic/gin"
	"sync"
)

const door = ":8080"

func HandleRequests() {
	r := gin.Default()
	r.POST("/name", controllers.CreateName)
	r.DELETE("/:id", controllers.DeleteName)
	r.PATCH("/:id", controllers.UpdateName)
	r.GET("/:id", waitGroupID)
	r.GET("/name/:name", waitGroupName)
	r.GET("/metaphone/:name", waitGroupMetaphone)

	err := r.Run(door)
	if err != nil {
		panic(err)
	}
}

//waitGroupMetaphone crates a waiting group for handling requests using controllers.SearchSimilarNames
func waitGroupMetaphone(c *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(1)

	// Handle the request in a separate goroutine
	go func() {
		defer wg.Done()
		controllers.SearchSimilarNames(c)
	}()

	wg.Wait()
}

//waitGroupMetaphone crates a waiting group for handling requests using controllers.GetName
func waitGroupName(c *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(1)

	// Handle the request in a separate goroutine
	go func() {
		defer wg.Done()
		controllers.GetName(c)
	}()

	wg.Wait()
}

//waitGroupMetaphone crates a waiting group for handling requests using controllers.GetID
func waitGroupID(c *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(1)

	// Handle the request in a separate goroutine
	go func() {
		defer wg.Done()
		controllers.GetID(c)
	}()

	wg.Wait()
}
