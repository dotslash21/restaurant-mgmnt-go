package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrderItems(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetOrderItems"})
}

func GetOrderItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetOrderItem"})
}

func GetOrderItemsByOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetOrderItemsByOrder"})
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	return nil, nil
}

func CreateOrderItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "CreateOrderItem"})
}

func UpdateOrderItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "UpdateOrderItem"})
}
