package main

import (
	"app/api"
	"app/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// format loggin
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// connect database
	db.ConnectDatabase()

	// set up router
	router := gin.New()

	// setup API
	api.Setup(router)

	log.Fatal(router.Run(":8080"))
}
