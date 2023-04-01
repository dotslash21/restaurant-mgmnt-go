package routers

import (
	"restaurant-mgmnt/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemRouter(router *gin.Engine) {
	router.GET("/orderItems", controllers.GetOrderItems)
	router.GET("/orderItems/:id", controllers.GetOrderItem)
	router.GET("/orderItems/order/:id", controllers.GetOrderItemsByOrder)
	router.POST("/orderItems", controllers.CreateOrderItem)
	router.PATCH("/orderItems/:id", controllers.UpdateOrderItem)
}
