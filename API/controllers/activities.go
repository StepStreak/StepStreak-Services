package controllers

import (
	"net/http"
	"strconv"
	"time"

	"example.com/main/messages"
	"example.com/main/models"

	"github.com/gin-gonic/gin"
)

type ActivityInput struct {
	Date  time.Time `json:"date" binding:"required"`
	Type  string    `json:"type" binding:"required"`
	Unit  string    `json:"unit" binding:"required"`
	Value string    `json:"value" binding:"required"`
}

type CreateActivitiesInput struct {
	Data []ActivityInput `json:"data" binding:"required"`
}

func CreateActivity(c *gin.Context) {
	var input CreateActivitiesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, activityInput := range input.Data {
		value, _ := strconv.ParseFloat(activityInput.Value, 64)

		var existingActivity models.Activity
		if err := models.DB.Where("date = ? AND type = ?", activityInput.Date, activityInput.Type).First(&existingActivity).Error; err != nil {
			activity := models.Activity{Date: activityInput.Date, Type: activityInput.Type, Unit: activityInput.Unit, Value: int32(value)}
			models.DB.Create(&activity)
			messages.Produce(activity.ID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": input.Data})
}
