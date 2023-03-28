package routes

import (
	"github.com/Darklabel91/API_Names/controllers"
	"github.com/Darklabel91/API_Names/middleware"
	"github.com/gin-gonic/gin"
	"sync"
)

const door = ":8080"

func HandleRequests() {
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/name", middleware.RequireAuth, controllers.CreateName)
	r.DELETE("/:id", middleware.RequireAuth, controllers.DeleteName)
	r.PATCH("/:id", middleware.RequireAuth, controllers.UpdateName)
	r.GET("/:id", middleware.RequireAuth, WaitGroupID)
	r.GET("/name/:name", middleware.RequireAuth, WaitGroupName)
	r.GET("/metaphone/:name", middleware.RequireAuth, WaitGroupMetaphone)

	err := r.Run(door)
	if err != nil {
		panic(err)
	}
}

//WaitGroupMetaphone crates a waiting group for handling requests using controllers.SearchSimilarNames
func WaitGroupMetaphone(c *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(1)

	// Handle the request in a separate goroutine
	go func() {
		defer wg.Done()
		controllers.SearchSimilarNames(c)
	}()

	wg.Wait()
}

//WaitGroupName crates a waiting group for handling requests using controllers.GetName
func WaitGroupName(c *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(1)

	// Handle the request in a separate goroutine
	go func() {
		defer wg.Done()
		controllers.GetName(c)
	}()

	wg.Wait()
}

// WaitGroupID  crates a waiting group for handling requests using controllers.GetID
func WaitGroupID(c *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(1)

	// Handle the request in a separate goroutine
	go func() {
		defer wg.Done()
		controllers.GetID(c)
	}()

	wg.Wait()
}
