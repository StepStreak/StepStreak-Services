package models

import "time"

type Activity struct {
	ID    uint      `json:"id" gorm:"primary_key"`
	Steps int16     `json:"steps"`
	Date  time.Time `json:"date"`
}
