package routes

import (
	"github.com/Darklabel91/API_Names/controllers"
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/middleware"
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"sync"
)

const door = ":8080"

var allowedIPs = []string{"127.0.0.1", "::1"} // List of allowed IP addresses

func HandleRequests() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Use the OnlyAllowIPs middleware on all routes
	err := r.SetTrustedProxies(allowedIPs)
	if err != nil {
		return
	}

	//signup and login
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	//Other routes
	r.Use(middleware.RequireAuth)

	r.POST("/name", controllers.CreateName)
	r.DELETE("/:id", controllers.DeleteName)
	r.PATCH("/:id", controllers.UpdateName)
	r.GET("/:id", WaitGroupID)
	r.GET("/name/:name", WaitGroupName)
	r.GET("/metaphone/:name", PreloadNameTypes(), WaitGroupMetaphone)

	err = r.Run(door)
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

//PreloadNameTypes for better response time we load all records of the table
func PreloadNameTypes() gin.HandlerFunc {
	var nameTypes []models.NameType
	if err := database.Db.Find(&nameTypes).Error; err != nil {
		return nil
	}

	return func(c *gin.Context) {
		c.Set("nameTypes", nameTypes)
		c.Next()
	}
}
