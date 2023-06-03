package main

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Firstname   string             `json:"firstname" bson:"firstname"`
	Lastname    string             `json:"lastname" bson:"lastname"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	Password    string             `json:"password" bson:"password"`
}

var Client *mongo.Client

func main() {
	// format loggin
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// database connect
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	Client = client

	ctx, cacnel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cacnel()
	err = Client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// set up router
	router := gin.New()

	router.POST("/new", func(c *gin.Context) {
		var user User
		err := c.Bind(&user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		usersColl := Client.Database("hala").Collection("users")

		result := usersColl.FindOne(context.TODO(), bson.M{
			"phone_number": user.PhoneNumber,
		})
		if result.Err() != nil {
			if result.Err() == mongo.ErrNoDocuments {
				// create the user
				user.ID = primitive.NewObjectID()
				user.Password = hashString(user.Password)

				res, err := usersColl.InsertOne(context.TODO(), user)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"user":         user,
					"insertion_id": res.InsertedID,
				})
				return
			}
			// internal server error
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		// user exists

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "user with this phone number already exists",
		})
	})

	log.Fatal(router.Run(":8080"))
}

func hashString(stringToHash string) string {
	salt := "salt is bad for health"
	hasher := sha1.New()
	hasher.Write([]byte(stringToHash + salt))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
