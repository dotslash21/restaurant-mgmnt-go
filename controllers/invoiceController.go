package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetInvoices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetInvoices"})
}

func GetInvoice(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetInvoice"})
}

func CreateInvoice(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "CreateInvoice"})
}

func UpdateInvoice(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "UpdateInvoice"})
}
