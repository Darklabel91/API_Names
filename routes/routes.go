package routes

import (
	"github.com/Darklabel91/API_Names/controllers"
	"github.com/Darklabel91/API_Names/middlewares"
	"github.com/Darklabel91/API_Names/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

const DOOR = ":8080"
const FILENAME = "Logs.txt"
const MICROSECONDS = 300

var nameTypesCache []models.NameType

func HandleRequests() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//use the OnlyAllowIPs middlewares on all routes
	err := r.SetTrustedProxies(models.IPs)
	if err != nil {
		return
	}

	// Create a file to store the logs
	file, err := os.OpenFile(FILENAME, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	r.Use(gin.LoggerWithWriter(file))

	//upload the log file from time to time
	var log models.Log
	ticker := time.NewTicker(MICROSECONDS * time.Microsecond)
	defer ticker.Stop()
	log.UploadLog(ticker, FILENAME)

	//routes without middlewares
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	//main middleware validation
	r.Use(middlewares.ValidateAuth())

	//CRUD routes
	r.POST("/name", middlewares.ValidateName(), middlewares.ValidateNameJSON(), controllers.CreateName)
	r.GET("/:id", middlewares.ValidateID(), controllers.GetID)
	r.GET("/name/:name", middlewares.ValidateName(), controllers.GetName)
	r.GET("/metaphone/:name", middlewares.ValidateName(), cachingNameTypes(nameTypesCache), controllers.GetMetaphoneMatch)
	r.PATCH("/:id", middlewares.ValidateID(), middlewares.ValidateNameJSON(), controllers.UpdateName)
	r.DELETE("/:id", middlewares.ValidateID(), controllers.DeleteName)

	// run
	err = r.Run(DOOR)
	if err != nil {
		return
	}
}

func cachingNameTypes(nameTypesCache []models.NameType) gin.HandlerFunc {
	var name models.NameType

	if nameTypesCache == nil {
		nameTypes, err := name.GetAllNames()
		if err != nil {
			return func(c *gin.Context) {
				c.JSON(http.StatusInternalServerError, gin.H{"Message": "Error on caching all name types"})
			}
		}
		nameTypesCache = nameTypes
	}

	return func(c *gin.Context) {
		c.Set("nameTypes", nameTypesCache)
		c.Next()
	}
}
