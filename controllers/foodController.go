package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetFoods(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetFoods"})
}

func GetFood(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetFood"})
}

func CreateFood(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "CreateFood"})
}

func round(num float64) int {
	return 0
}

func toFixed(num float64, precision int) float64 {
	return 0.0
}

func UpdateFood(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "UpdateFood"})
}
