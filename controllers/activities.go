package controllers

import (
	"net/http"
	"time"

	"example.com/main/messages"
	"example.com/main/models"

	"github.com/gin-gonic/gin"
)

type CreateActivityInput struct {
	Date  time.Time `json:"date" binding:"required"`
	Steps int16     `json:"steps" binding:"required"`
}

func CreateActivity(c *gin.Context) {
	var input CreateActivityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activity := models.Activity{Date: input.Date, Steps: input.Steps}
	models.DB.Create(&activity)

	messages.Produce(activity.ID)

	c.JSON(http.StatusOK, gin.H{"data": activity})
}
