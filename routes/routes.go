package routes

import (
	"github.com/Darklabel91/API_Names/controllers"
	"github.com/gin-gonic/gin"
)

const door = ":8080"

func HandleRequests() {
	r := gin.Default()
	r.POST("/name", controllers.CreateName)
	r.DELETE("/:id", controllers.DeleteName)
	r.PATCH("/:id", controllers.UpdateName)
	r.GET("/:id", controllers.SearchNameByID)
	r.GET("/name/:name", controllers.GetName)
	r.GET("/metaphone/:name", controllers.SearchSimilarNames)

	err := r.Run(door)
	if err != nil {
		panic(err)
	}
}
