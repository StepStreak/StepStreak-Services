package messages

import (
	"encoding/json"
	"time"

	"example.com/main/models"

	"github.com/streadway/amqp"
)

type Activity struct {
	Steps int16     `json:"steps"`
	Date  time.Time `json:"date"`
}

func Produce(id uint) {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	q, _ := ch.QueueDeclare(
		"activities",
		false,
		false,
		false,
		false,
		nil,
	)

	var activity models.Activity

	models.DB.Where("id = ?", id).First(&activity)

	body, _ := json.Marshal(activity)

	_ = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
