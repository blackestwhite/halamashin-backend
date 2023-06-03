package handlers

import (
	"app/entity"
	"app/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userService service.UserService
}

func SetupUser(r *gin.RouterGroup) *UserHandler {
	u := &UserHandler{}
	u.initRoutes(r)
	return u
}

func (u *UserHandler) initRoutes(r *gin.RouterGroup) {
	user := r.Group("user")

	user.POST("/new", u.createUser)
	user.GET("/get/:phone_number", u.getByPhoneNumber)
}

func (u *UserHandler) getByPhoneNumber(c *gin.Context) {
	user, err := u.userService.GetByPhoneNumber(c.Param("phone_number"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (u *UserHandler) createUser(c *gin.Context) {
	var user entity.User
	err := c.Bind(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	_, err = u.userService.GetByPhoneNumber(user.PhoneNumber)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			// internal server error
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		// create the user
		insertionID, err := u.userService.Create(user)
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
}
