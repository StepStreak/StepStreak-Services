package main

import (
	"net/http"

	"example.com/main/controllers"
	"example.com/main/models"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	r.POST("/activities", controllers.CreateActivity)

	models.ConnectDatabase()

	r.Run()
}
