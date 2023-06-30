package models

import "time"

type Activity struct {
	ID    uint      `json:"id" gorm:"primary_key"`
	Type  string    `json:"type"`
	Unit  string    `json:"unit"`
	Value int32     `json:"value"`
	Date  time.Time `json:"date"`
}
