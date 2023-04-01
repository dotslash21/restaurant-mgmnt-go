package routers

import (
	"restaurant-mgmnt/controllers"

	"github.com/gin-gonic/gin"
)

func InvoiceRouter(router *gin.Engine) {
	router.GET("/invoices", controllers.GetInvoices)
	router.GET("/invoices/:id", controllers.GetInvoice)
	router.POST("/invoices", controllers.CreateInvoice)
	router.PATCH("/invoices/:id", controllers.UpdateInvoice)
}
