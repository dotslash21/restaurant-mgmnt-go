package controllers

import (
	"context"
	"log"
	"net/http"
	"restaurant-mgmnt/helpers"
	"restaurant-mgmnt/models"
	"restaurant-mgmnt/repositories"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = repositories.OpenCollection(repositories.Client, "user")

func GetUsers(c *gin.Context) {
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

	if result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage}); err != nil {
		c.JSON(500, gin.H{
			"message": "Error occurred while fetching the user items.",
		})
	} else {
		var users []bson.M
		if err := result.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error occurred while fetching the user items.",
			})
		} else {
			c.JSON(http.StatusOK, users)
		}
	}
}

func GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	id := c.Param("id")

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Error occurred while fetching the user item.",
		})
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func SignUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// 1. Parse the sign up data from the request
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error occurred while parsing the sign up data.",
		})
		return
	}

	// 2. Validate the data
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 3. If email already exists, return error
	if count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occurred while fetching the user items.",
		})
		return
	} else if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already exists.",
		})
		return
	}

	// 4. Hash the password
	password := HashPassword(*user.Password)
	user.Password = &password

	// 5. Check if the phone number already exists, return error
	if count, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occurred while fetching the user items.",
		})
		return
	} else if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Phone number already exists.",
		})
		return
	}

	// 6. Fill out the audit fields - created_at, updated_at
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Id = primitive.NewObjectID()

	// 7. Generate tokens
	token, refreshToken, err := helpers.GenerateTokens(*user.Email, *user.FirstName, *user.LastName, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occurred while generating the tokens.",
		})
		return
	}
	user.Token = &token
	user.RefreshToken = &refreshToken

	// 8. Insert the user into the database
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occurred while inserting the user item.",
		})
		return
	}

	// 9. Return the saved user item with status code 201
	c.JSON(http.StatusCreated, result)
}

func LogIn(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// 1. Parse the login data from the request
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error occurred while parsing the login data.",
		})
		return
	}

	// 2. Find the user by email
	var foundUser models.User
	if err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Error occurred while fetching the user item.",
		})
		return
	}

	// 4. Verify the password
	isValidPassword, message := VerifyPassword(*foundUser.Password, *user.Password)
	if !isValidPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	// 5. Generate tokens
	token, refreshToken, err := helpers.GenerateTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occurred while generating the tokens.",
		})
		return
	}

	// 6. Update the user with the new tokens
	helpers.UpdateTokens(token, refreshToken, foundUser.Id)

	// 7. Return the tokens with status code 200
	c.JSON(http.StatusOK, foundUser)
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(hashedPassword, plainTextPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))
	if err != nil {
		return false, "Invalid password."
	}

	return true, "Password is valid."
}
