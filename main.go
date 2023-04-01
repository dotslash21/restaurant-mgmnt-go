package main

import (
	"os"
	"restaurant-mgmnt/middlewares"
	"restaurant-mgmnt/repositories"
	"restaurant-mgmnt/routers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = repositories.OpenCollection(repositories.Client, "food")

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routers.UserRouter(router)

	// Protected routes - start
	router.Use(middlewares.Authentication())

	routers.FoodRouter(router)
	routers.InvoiceRouter(router)
	routers.MenuRouter(router)
	routers.OrderRouter(router)
	routers.OrderItemRouter(router)
	routers.TableRouter(router)
	// Protected routes - end

	router.Run(":" + port)
}
