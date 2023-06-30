package messages

import (
	"encoding/json"
	"log"
	"time"

	"example.com/main/models"

	"github.com/streadway/amqp"
)

type Activity struct {
	Type  string    `json:"type"`
	Unit  string    `json:"unit"`
	Value int32     `json:"value"`
	Date  time.Time `json:"date"`
}

func Produce(id uint) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Printf("Failed to connect to the RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"messages",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return
	}

	var activity models.Activity
	models.DB.Where("id = ?", id).First(&activity)

	body, err := json.Marshal(activity)
	if err != nil {
		log.Printf("Failed to marshal the activity: %v", err)
		return
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
		return
	}
}
