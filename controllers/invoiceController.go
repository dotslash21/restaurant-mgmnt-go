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

var invoiceCollection *mongo.Collection = repositories.OpenCollection(repositories.Client, "invoice")

func GetInvoices(c *gin.Context) {
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

	if result, err := invoiceCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		groupStage,
		projectStage,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the invoice items."})
	} else {
		var invoices []bson.M
		if err := result.All(ctx, &invoices); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the invoice items."})
		}
		c.JSON(http.StatusOK, invoices)
	}
}

func GetInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error occurred while parsing the invoice id.",
		})
		return
	}

	var invoice models.Invoice
	err = invoiceCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&invoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the invoice item."})
	}

	var invoiceView models.InvoiceView
	invoiceView.Id = invoice.Id
	invoiceView.OrderId = invoice.OrderId
	invoiceView.PaymentDueDate = invoice.PaymentDueDate
	if invoice.PaymentMethod != nil {
		invoiceView.PaymentMethod = *invoice.PaymentMethod
	} else {
		invoiceView.PaymentMethod = "null"
	}
	invoiceView.PaymentStatus = invoice.PaymentStatus
	if orderItems, err := ItemsByOrder(ctx, invoice.OrderId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the invoice item."})
	} else {
		invoiceView.PaymentDue = orderItems[0]["payment_due"]
		invoiceView.TableNumber = orderItems[0]["table_number"]
		invoiceView.OrderDetails = orderItems[0]["order_items"]
	}

	c.JSON(http.StatusOK, invoiceView)
}

func CreateInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var invoice models.Invoice
	if err := c.ShouldBindJSON(&invoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while binding the invoice item."})
		return
	}

	var order models.Order
	if err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.OrderId}).Decode(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the order item."})
	}

	invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.Id = primitive.NewObjectID()
	paymentStatus := "PENDING"
	if invoice.PaymentStatus == nil {
		invoice.PaymentStatus = &paymentStatus
	}
	invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))

	if err := validate.Struct(invoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while validating the invoice item."})
		return
	}

	if result, err := invoiceCollection.InsertOne(ctx, invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while inserting the invoice item."})
		return
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func UpdateInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var invoice models.Invoice
	if err := c.BindJSON(&invoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while binding the invoice item."})
		return
	}

	var updateObj primitive.D
	if invoice.PaymentMethod != nil {
		updateObj = append(updateObj, bson.E{Key: "payment_method", Value: *invoice.PaymentMethod})
	}
	if invoice.PaymentStatus != nil {
		updateObj = append(updateObj, bson.E{Key: "payment_status", Value: *invoice.PaymentStatus})
	} else {
		updateObj = append(updateObj, bson.E{Key: "payment_status", Value: "PENDING"})
	}
	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: invoice.UpdatedAt})

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error occurred while parsing the invoice id.",
		})
		return
	}

	filter := bson.M{"_id": id}
	upsert := true
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}
	if result, err := invoiceCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating the invoice item."})
		return
	} else {
		c.JSON(http.StatusOK, result)
	}
}
