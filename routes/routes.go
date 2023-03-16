package routes

import (
	"github.com/Darklabel91/API_Names/controllers"
	"github.com/gin-gonic/gin"
)

const door = ":8080"

func HandleRequests() {
	r := gin.Default()
	r.GET("apiName/", controllers.GetAllNames)

	r.POST("apiName/name", controllers.CreateName)
	r.GET("apiName/:id", controllers.SearchNameByID)
	r.DELETE("apiName/:id", controllers.DeleteName)
	r.PATCH("apiName/:id", controllers.UpdateName)
	r.GET("apiName/metaphone/:mtf", controllers.SearchNameByMetaphone)

	err := r.Run(door)
	if err != nil {
		panic(err)
	}
}
