package routers

import (
	"restaurant-mgmnt/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRouter(router *gin.Engine) {
	router.GET("/orders", controllers.GetOrders)
	router.GET("/orders/:id", controllers.GetOrder)
	router.POST("/orders", controllers.CreateOrder)
	router.PATCH("/orders/:id", controllers.UpdateOrder)
}
