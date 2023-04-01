package routers

import (
	"restaurant-mgmnt/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRouter(router *gin.Engine) {
	router.GET("/menus", controllers.GetMenus)
	router.GET("/menus/:id", controllers.GetMenu)
	router.POST("/menus", controllers.CreateMenu)
	router.PATCH("/menus/:id", controllers.UpdateMenu)
}
