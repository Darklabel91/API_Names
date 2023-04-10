package routes

import (
	"github.com/Darklabel91/API_Names/controllers"
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/middleware"
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
	"time"
)

const DOOR = ":8080"
const FILENAME = "Logs.txt"
const MINUTES = 1

var nameTypesCache []models.NameType

func HandleRequests() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//use the OnlyAllowIPs middleware on all routes
	err := r.SetTrustedProxies(controllers.GetTrustedIPs())
	if err != nil {
		return
	}

	// Create a file to store the logs
	file, err := os.OpenFile(FILENAME, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	r.Use(gin.LoggerWithWriter(file))

	//upload the log file every given number of minutes
	var log models.Log
	ticker := time.NewTicker(MINUTES * time.Minute)
	defer ticker.Stop()
	uploadLog(log, database.DB, ticker, FILENAME)

	//set up routes
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	//define middleware that validate token
	r.Use(middleware.ValidateAuth())

	//set up caching middleware for GET requests
	r.DELETE("/:id", middleware.ValidateID(), controllers.DeleteName)
	r.PATCH("/:id", middleware.ValidateID(), controllers.UpdateName)
	r.POST("/name", middleware.ValidateNameJSON(), controllers.CreateName)
	r.GET("/:id", middleware.ValidateID(), controllers.GetID)
	r.GET("/name/:name", middleware.ValidateName(), controllers.GetName)
	r.GET("/metaphone/:name", middleware.ValidateName(), cachingNameTypes(), controllers.GetSimilarNames)

	// run
	err = r.Run(DOOR)
	if err != nil {
		return
	}
}

//uploadLog the logger
func uploadLog(log models.Log, db *gorm.DB, ticker *time.Ticker, fileName string) {
	go func() {
		for {
			select {
			case <-ticker.C:
				err := log.Upload(db, fileName)
				if err != nil {
					panic(err)
				}

			}
		}
	}()
}

//cachingNameTypes for better response time we load all records of the table
func cachingNameTypes() gin.HandlerFunc {
	if nameTypesCache == nil {
		var nameTypes []models.NameType
		nameTypesCache = nameTypes

		if err := database.DB.Find(&nameTypes).Error; err != nil {
			return nil
		}

	}

	return func(c *gin.Context) {
		c.Set("nameTypes", nameTypesCache)
		c.Next()
	}
}
