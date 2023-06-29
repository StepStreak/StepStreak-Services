package main

import (
	"example.com/main/controllers"
	"example.com/main/models"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/activities", controllers.CreateActivity)

	models.ConnectDatabase()

	r.Run()
}
