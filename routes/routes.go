package routes

import (
	"github.com/Darklabel91/API_Names/controllers"
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/middleware"
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"os"
	"sync"
)

const DOOR = ":8080"
const FILENAME = "logs.txt"

var allowedIPs = []string{"127.0.0.1", "::1"} // List of allowed IP addresses

func HandleRequests() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//use the OnlyAllowIPs middleware on all routes
	err := r.SetTrustedProxies(allowedIPs)
	if err != nil {
		return
	}

	// Create a file to store the logs
	file, err := os.OpenFile(FILENAME, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	r.Use(gin.LoggerWithWriter(file))

	//set up routes
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	//define middleware that validate token
	r.Use(middleware.RequireAuth())

	//set up caching middleware for GET requests
	r.GET("/:id", middleware.ValidateIDParam(), waitGroupID)
	r.DELETE("/:id", middleware.ValidateIDParam(), controllers.DeleteName)
	r.PATCH("/:id", middleware.ValidateIDParam(), controllers.UpdateName)
	r.POST("/name", controllers.CreateName)
	r.GET("/name/:name", middleware.ValidateNameParam(), waitGroupName)
	r.GET("/metaphone/:name", middleware.ValidateNameParam(), preloadNameTypes(), middleware.ValidateNameParam(), waitGroupMetaphone)

	gin.Logger()

	// run
	err = r.Run(DOOR)
	if err != nil {
		return
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

//waitGroupName crates a waiting group for handling requests using controllers.GetName
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

// waitGroupID  crates a waiting group for handling requests using controllers.GetID
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

//preloadNameTypes for better response time we load all records of the table
func preloadNameTypes() gin.HandlerFunc {
	var nameTypes []models.NameType
	if err := database.Db.Find(&nameTypes).Error; err != nil {
		return nil
	}

	return func(c *gin.Context) {
		c.Set("nameTypes", nameTypes)
		c.Next()
	}
}
