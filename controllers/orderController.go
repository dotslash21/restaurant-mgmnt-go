package controllers

import (
	"context"
	"net/http"
	"restaurant-mgmnt/models"
	"restaurant-mgmnt/repositories"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection *mongo.Collection = repositories.OpenCollection(repositories.Client, "order")

func GetOrders(c *gin.Context) {
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

	if result, err := orderCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		groupStage,
		projectStage,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the order items."})
	} else {
		var orders []bson.M
		if err := result.All(ctx, &orders); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the order items."})
		}
		c.JSON(http.StatusOK, orders)
	}
}

func GetOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Param("id")

	var order models.Order
	err := foodCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the order item."})
	}

	c.JSON(http.StatusOK, order)
}

func CreateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while binding the order item."})
		return
	}

	var table models.Table
	if order.TableId != nil {
		if err := tableCollection.FindOne(ctx, bson.M{"_id": order.TableId}).Decode(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the table."})
			return
		}
	}

	if err := validate.Struct(order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while validating the order item."})
		return
	}

	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Id = primitive.NewObjectID()

	if result, err := orderCollection.InsertOne(ctx, order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while creating the order item."})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func UpdateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while binding the order item."})
		return
	}

	var updateObj primitive.D
	if order.TableId != nil {
		var table models.Table
		if err := tableCollection.FindOne(ctx, bson.M{"_id": order.TableId}).Decode(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the table."})
			return
		}
		updateObj = append(updateObj, bson.E{Key: "table_id", Value: order.TableId})
	}
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: order.UpdatedAt})

	id := c.Param("id")
	filter := bson.M{"_id": id}
	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}
	if result, err := orderCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating the order item."})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func OrderItemOrderCreation(order models.Order) primitive.ObjectID {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Id = primitive.NewObjectID()

	if result, err := orderCollection.InsertOne(ctx, order); err != nil {
		return primitive.NilObjectID
	} else {
		return result.InsertedID.(primitive.ObjectID)
	}
}
