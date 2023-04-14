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

var orderItemCollection *mongo.Collection = repositories.OpenCollection(repositories.Client, "order_item")

func GetOrderItems(c *gin.Context) {
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

	if result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		groupStage,
		projectStage,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the orderItem items."})
	} else {
		var orderItems []bson.M
		if err := result.All(ctx, &orderItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the orderItem items."})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": orderItems})
		}
	}
}

func GetOrderItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Param("id")

	var orderItem models.OrderItem
	err := orderItemCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&orderItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the orderItem item."})
	}

	c.JSON(http.StatusOK, orderItem)
}

func GetOrderItemsByOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while parsing the order ID."})
	}

	if orderItems, err := ItemsByOrder(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the orderItems by order ID."})
	} else {
		c.JSON(http.StatusOK, orderItems)
	}
}

func ItemsByOrder(ctx context.Context, id primitive.ObjectID) (orderItems []primitive.M, err error) {
	matchStage := bson.D{
		{
			Key: "$match",
			Value: bson.D{
				{
					Key:   "_id",
					Value: id,
				},
			},
		},
	}

	foodLookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{
					Key:   "from",
					Value: "food",
				},
				{
					Key:   "localField",
					Value: "food_id",
				},
				{
					Key:   "foreignField",
					Value: "_id",
				},
				{
					Key:   "as",
					Value: "food",
				},
			},
		},
	}
	foodUnwindStage := bson.D{
		{
			Key: "$unwind",
			Value: bson.D{
				{
					Key:   "path",
					Value: "$food",
				},
				{
					Key:   "preserveNullAndEmptyArrays",
					Value: true,
				},
			},
		},
	}

	orderLookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "order"},
				{Key: "localField", Value: "order_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "order"},
			},
		},
	}
	orderUnwindStage := bson.D{
		{
			Key: "$unwind",
			Value: bson.D{
				{Key: "path", Value: "$order"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			},
		},
	}

	tableLookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "table"},
				{Key: "localField", Value: "table_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "table"},
			},
		},
	}
	tableUnwindStage := bson.D{
		{
			Key: "$unwind",
			Value: bson.D{
				{Key: "path", Value: "$table"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			},
		},
	}

	projectStage := bson.D{
		{
			Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "amount", Value: "$food.price"},
				{Key: "total_count", Value: 1},
				{Key: "food_name", Value: "$food.name"},
				{Key: "food_image", Value: "$food.image"},
				{Key: "table_number", Value: "$table.number"},
				{Key: "table_id", Value: "$table._id"},
				{Key: "order_id", Value: "$order._id"},
				{Key: "price", Value: "$food.price"},
				{Key: "quantity", Value: 1},
			},
		},
	}

	groupStage := bson.D{
		{
			Key: "$group",
			Value: bson.D{
				{
					Key: "_id",
					Value: bson.D{
						{Key: "order_id", Value: "$order_id"},
						{Key: "table_id", Value: "$table_id"},
						{Key: "table_number", Value: "$table_number"},
						{Key: "payment_due", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
						{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
						{Key: "order_items", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
					},
				},
			},
		},
	}

	projectStage2 := bson.D{
		{
			Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "payment_due", Value: 1},
				{Key: "total_count", Value: 1},
				{Key: "table_number", Value: "$_id.table_number"},
				{Key: "order_items", Value: 1},
			},
		},
	}

	if result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		foodLookupStage,
		foodUnwindStage,
		orderLookupStage,
		orderUnwindStage,
		tableLookupStage,
		tableUnwindStage,
		projectStage,
		groupStage,
		projectStage2,
	}); err != nil {
		panic(err)
	} else {
		if err = result.All(ctx, &orderItems); err != nil {
			panic(err)
		}
	}

	return orderItems, err
}

func CreateOrderItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var orderItem models.OrderItem
	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while binding the orderItem item."})
		return
	}

	if err := validate.Struct(orderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while validating the orderItem item."})
		return
	}

	orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderItem.Id = primitive.NewObjectID()

	if result, err := orderItemCollection.InsertOne(ctx, orderItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while creating the orderItem item."})
	} else {
		c.JSON(http.StatusCreated, result)
	}
}

func UpdateOrderItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var orderItem models.OrderItem
	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while binding the orderItem item."})
		return
	}

	var updateObj primitive.D
	if orderItem.UnitPrice != nil {
		updateObj = append(updateObj, bson.E{Key: "unit_price", Value: orderItem.UnitPrice})
	}
	if orderItem.Quantity != nil {
		updateObj = append(updateObj, bson.E{Key: "quantity", Value: orderItem.Quantity})
	}
	if orderItem.FoodId != nil {
		updateObj = append(updateObj, bson.E{Key: "food_id", Value: orderItem.FoodId})
	}
	orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "$set", Value: orderItem})

	id := c.Param("id")
	filter := bson.M{"_id": id}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	if result, err := orderItemCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating the orderItem item."})
	} else {
		c.JSON(http.StatusOK, result)
	}
}
