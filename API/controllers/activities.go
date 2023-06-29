package controllers

import (
	"net/http"
	"strconv"
	"time"

	"example.com/main/messages"
	"example.com/main/models"

	"github.com/gin-gonic/gin"
)

type CreateActivityInput struct {
	Date  time.Time `json:"date" binding:"required"`
	Type  string    `json:"type" binding:"required"`
	Unit  string    `json:"unit" binding:"required"`
	Value string    `json:"value" binding:"required"`
}

func CreateActivity(c *gin.Context) {
	var input CreateActivityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value, _ := strconv.ParseInt(input.Value, 10, 16)

	activity := models.Activity{Date: input.Date, Type: input.Type, Unit: input.Unit, Value: int16(value)}
	models.DB.Create(&activity)

	messages.Produce(activity.ID)

	c.JSON(http.StatusOK, gin.H{"data": activity})
}
