package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type IncomingMessage struct {
	ID    int       `json:"id"`
	Type  string    `json:"type"`
	Unit  string    `json:"unit"`
	Value int       `json:"value"`
	Date  time.Time `json:"date"`
}

type CombinedMessage struct {
	Steps          int    `json:"steps"`
	ActiveCalories int    `json:"active_calories"`
	Date           string `json:"date"`
}

var messageMap = make(map[string]map[string]IncomingMessage)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"messages", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	outgoingQ, err := ch.QueueDeclare(
		"activities", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		for d := range msgs {
			var msg IncomingMessage
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				log.Printf("Error decoding message: %v", err)
				continue
			}

			dateStr := msg.Date.Format("2006-01-02")
			if _, ok := messageMap[dateStr]; !ok {
				messageMap[dateStr] = make(map[string]IncomingMessage)
			}

			messageMap[dateStr][msg.Type] = msg

			if len(messageMap[dateStr]) == 2 {
				stepsMessage, stepsOk := messageMap[dateStr]["Steps"]
				caloriesMessage, caloriesOk := messageMap[dateStr]["Active Calories"]
				if stepsOk && caloriesOk {
					combinedMessage := CombinedMessage{
						Steps:          stepsMessage.Value,
						ActiveCalories: caloriesMessage.Value,
						Date:           dateStr,
					}

					body, err := json.Marshal(combinedMessage)
					if err != nil {
						log.Printf("Error encoding message: %v", err)
						continue
					}

					err = ch.Publish(
						"",             // exchange
						outgoingQ.Name, // routing key
						false,          // mandatory
						false,          // immediate
						amqp.Publishing{
							ContentType: "application/json",
							Body:        body,
						})
					if err != nil {
						log.Printf("Failed to publish a message: %v", err)
					}

					d.Ack(false)
					// Assuming you stored delivery tag for the first message, you should ack it as well.
					// For example: stepsMessage.DeliveryTag.Ack(false)

					delete(messageMap, dateStr)
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	forever := make(chan bool)
	<-forever
}
