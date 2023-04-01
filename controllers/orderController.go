package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetOrders"})
}

func GetOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetOrder"})
}

func CreateOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "CreateOrder"})
}

func UpdateOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "UpdateOrder"})
}
