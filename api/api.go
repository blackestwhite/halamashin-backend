package api

import (
	"app/api/handlers"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {
	apiGroup := router.Group("api")
	v1 := apiGroup.Group("v1")

	handlers.SetupStats(v1)
	handlers.SetupUser(v1)
}
