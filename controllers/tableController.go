package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTables(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetTables"})
}

func GetTable(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetTable"})
}

func CreateTable(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "CreateTable"})
}

func UpdateTable(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "UpdateTable"})
}
