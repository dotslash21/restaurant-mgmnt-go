package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMenus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetMenus"})
}

func GetMenu(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetMenu"})
}

func CreateMenu(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "CreateMenu"})
}

func UpdateMenu(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "UpdateMenu"})
}
