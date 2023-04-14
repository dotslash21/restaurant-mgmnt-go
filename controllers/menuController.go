package controllers

import (
	"context"
	"net/http"
	"restaurant-mgmnt/models"
	"restaurant-mgmnt/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection *mongo.Collection = repositories.OpenCollection(repositories.Client, "menu")

func GetMenus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if result, err := menuCollection.Find(ctx, bson.M{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the menus."})
	} else {
		var menus []bson.M
		if err := result.All(ctx, &menus); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the menus."})
		}
		c.JSON(http.StatusOK, menus)
	}
}

func GetMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error occurred while parsing the menu id.",
		})
		return
	}

	var menu models.Menu
	err = menuCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&menu)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the menu item."})
	}
	c.JSON(http.StatusOK, menu)
}

func CreateMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var menu models.Menu
	if err := c.BindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while binding the menu item."})
		return
	}

	if err := validate.Struct(menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while validating the menu item."})
		return
	}

	menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.Id = primitive.NewObjectID()

	if result, err := menuCollection.InsertOne(ctx, menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while inserting the menu item."})
		return
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func withInTimeSpan(start, end, check time.Time) bool {
	return start.After(check) && end.After(start)
}

func UpdateMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var menu models.Menu
	if err := c.BindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while binding the menu item."})
		return
	}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error occurred while parsing the menu id.",
		})
		return
	}

	filter := bson.M{"_id": id}

	var updateObj primitive.D
	if menu.StartDate != nil && menu.EndDate != nil {
		if !withInTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range for \"StartDate\" and \"EndDate\"."})
			return
		}

		updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.StartDate})
		updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.EndDate})

		if menu.Name != "" {
			updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
		}

		if menu.Category != "" {
			updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
		}

		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		if result, err := menuCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating the menu item."})
			return
		} else {
			c.JSON(http.StatusOK, result)
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range for \"StartDate\" and \"EndDate\"."})
	}
}
