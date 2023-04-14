package controllers

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"restaurant-mgmnt/models"
	"restaurant-mgmnt/repositories"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = repositories.OpenCollection(repositories.Client, "food")
var validate = validator.New()

func GetFoods(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	var startIndex int
	if value, err := strconv.Atoi(c.Query("startIndex")); err != nil {
		startIndex = (page - 1) * pageSize
	} else {
		startIndex = value
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{{
		Key: "$group",
		Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		},
	}}
	projectStage := bson.D{
		{
			Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "count", Value: 1},
				{
					Key: "data",
					Value: bson.D{
						{
							Key: "$slice",
							Value: bson.A{
								"$data",
								startIndex,
								pageSize,
							},
						},
					},
				},
			},
		},
	}

	if result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		groupStage,
		projectStage,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the food items."})
	} else {
		var foods []bson.M
		if err := result.All(ctx, &foods); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the food items."})
		}
		c.JSON(http.StatusOK, foods)
	}
}

func GetFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error occurred while parsing the food id.",
		})
		return
	}

	var food models.Food
	err = foodCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the food item."})
	}

	c.JSON(http.StatusOK, food)
}

func CreateFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var food models.Food
	if err := c.BindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while binding the food item."})
		return
	}

	var menu models.Menu
	if err := menuCollection.FindOne(ctx, bson.M{"_id": food.MenuId}).Decode(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while fetching the menu."})
		return
	}

	if err := validate.Struct(food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while validating the food item."})
		return
	}

	food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.Id = primitive.NewObjectID()
	price := toFixed(*food.Price, 2)
	food.Price = &price

	if result, err := foodCollection.InsertOne(ctx, food); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while inserting the food item."})
		return
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	rounded := float64(round(num*shift)) / shift
	return rounded
}

func UpdateFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var food models.Food
	if err := c.BindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while binding the food item."})
		return
	}

	var updateObj primitive.D
	if food.Name != nil {
		updateObj = append(updateObj, bson.E{Key: "name", Value: food.Name})
	}
	if food.Price != nil {
		updateObj = append(updateObj, bson.E{Key: "price", Value: food.Price})
	}
	if food.Image != nil {
		updateObj = append(updateObj, bson.E{Key: "image", Value: food.Image})
	}
	if food.MenuId != nil {
		var menu models.Menu
		if err := menuCollection.FindOne(ctx, bson.M{"_id": food.MenuId}).Decode(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while fetching the menu item."})
			return
		}

		updateObj = append(updateObj, bson.E{Key: "menu_id", Value: food.MenuId})
	}
	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: food.UpdatedAt})

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error occurred while parsing the food id.",
		})
		return
	}

	filter := bson.M{"_id": id}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	if result, err := foodCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating the food item."})
		return
	} else {
		c.JSON(http.StatusOK, result)
	}
}
