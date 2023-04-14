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

var tableCollection *mongo.Collection = repositories.OpenCollection(repositories.Client, "table")

func GetTables(c *gin.Context) {
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

	if result, err := tableCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		groupStage,
		projectStage,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the table items."})
	} else {
		var tables []interface{}
		if err := result.All(ctx, &tables); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the table items."})
		} else {
			c.JSON(http.StatusOK, tables)
		}
	}
}

func GetTable(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Param("id")

	var table models.Table
	err := tableCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&table)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}

	c.JSON(http.StatusOK, table)
}

func CreateTable(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var table models.Table
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while binding the table item."})
		return
	}

	if err := validate.Struct(table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while validating the table item."})
		return
	}

	table.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.Id = primitive.NewObjectID()

	if result, err := tableCollection.InsertOne(ctx, table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while creating the table item."})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func UpdateTable(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var table models.Table
	if err := c.BindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while binding the table item."})
		return
	}

	var updateObj primitive.D
	if table.NumberOfGuests != nil {
		updateObj = append(updateObj, bson.E{Key: "number_of_guests", Value: table.NumberOfGuests})
	}
	if table.Number != nil {
		updateObj = append(updateObj, bson.E{Key: "table_number", Value: table.Number})
	}
	table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: table.UpdatedAt})

	id := c.Param("id")
	filter := bson.M{"_id": id}
	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}
	if result, err := tableCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating the table item."})
	} else {
		c.JSON(http.StatusOK, result)
	}
}
