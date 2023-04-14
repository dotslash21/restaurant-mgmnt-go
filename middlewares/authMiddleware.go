package middlewares

import (
	"net/http"
	"restaurant-mgmnt/helpers"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "No Authorization header provided"})
		return
	}

	if !strings.HasPrefix(token, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header provided"})
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	claims, err := helpers.ValidateToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.Set("email", claims.Email)
	c.Set("firstName", claims.FirstName)
	c.Set("lastName", claims.LastName)
	c.Set("uid", claims.Uid)

	c.Next()
}
