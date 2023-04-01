package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetUsers"})
}

func GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "GetUser"})
}

func SignUp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "CreateUser"})
}

func LogIn(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "UpdateUser"})
}

func HashPassword(password string) string {
	return ""
}

func VerifyPassword(hashedPassword, password string) (bool, string) {
	return false, ""
}
