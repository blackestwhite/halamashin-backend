package handlers

import (
	"app/api/middleware"
	"app/entity"
	"app/service"
	"app/utils"
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

	user.POST("/create", u.createUser)
	user.GET("/get/:phone_number", u.getByPhoneNumber)
	user.PUT("/update/:user_id", middleware.Auth(), u.updateUser)
	user.POST("/login", u.loginUser)
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
			"message":      "ok, user created successfully",
			"insertion_id": insertionID,
		})
		return
	}

	// user exists
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"message": "user with this phone number already exists",
	})
}

func (u *UserHandler) updateUser(c *gin.Context) {
	fetchedUser, err := u.userService.Get(c.Param("user_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	value, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "you cannot access this resource",
		})
		return
	}

	if value.(string) != fetchedUser.ID.Hex() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "you cannot access this resource",
		})
		return
	}

	var user entity.User
	err = c.Bind(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	fetchedUser.Firstname = user.Firstname
	fetchedUser.Lastname = user.Lastname
	if utils.HashString(user.Password) != fetchedUser.Password && user.Password != "" {
		fetchedUser.Password = utils.HashString(user.Password)
	}

	err = u.userService.Update(fetchedUser, fetchedUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
	})
}

func (u *UserHandler) loginUser(c *gin.Context) {
	var user entity.User
	err := c.Bind(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	fetchedUser, err := u.userService.GetByPhoneNumber(user.PhoneNumber)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if utils.HashString(user.Password) != fetchedUser.Password {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "password mismatch",
		})
		return
	}

	tokenString, err := fetchedUser.GenerateToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}
	maxAge := 60 * 60 * 24 * 31

	// logged in successfully
	c.SetCookie(
		"hm-token",
		tokenString,
		maxAge,
		"/",
		"",
		false,
		false,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "logged in successfully",
	})
}
