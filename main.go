package main

import (
	"app/db"
	"app/entity"
	"app/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// format loggin
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// connect database
	db.ConnectDatabase()

	// set up router
	router := gin.New()

	router.POST("/new", func(c *gin.Context) {
		var user entity.User
		err := c.Bind(&user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		var userService service.UserService

		_, err = userService.GetByPhoneNumber(user.PhoneNumber)
		if err != nil {
			if err != mongo.ErrNoDocuments {
				// internal server error
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				return
			}
			// create the user
			insertionID, err := userService.Create(user)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"user":         user,
				"insertion_id": insertionID,
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
