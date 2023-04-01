package routers

import (
	"restaurant-mgmnt/controllers"

	"github.com/gin-gonic/gin"
)

func TableRouter(router *gin.Engine) {
	router.GET("/tables", controllers.GetTables)
	router.GET("/tables/:id", controllers.GetTable)
	router.POST("/tables", controllers.CreateTable)
	router.PATCH("/tables/:id", controllers.UpdateTable)
}
