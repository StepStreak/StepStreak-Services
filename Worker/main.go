package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type IncomingMessage struct {
	ID    int
	Type  string
	Unit  string
	Value int
	Date  time.Time
}

type OutgoingMessage struct {
	Steps          int
	ActiveCalories int
	Date           string
}

func main() {
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"messages", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	outgoingQ, err := ch.QueueDeclare(
		"activities", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	messageMap := make(map[string]map[string]amqp.Delivery)

	for d := range msgs {
		var msg IncomingMessage
		json.Unmarshal(d.Body, &msg)
		date := msg.Date.Format("2006-01-02")
		if _, ok := messageMap[date]; !ok {
			messageMap[date] = make(map[string]amqp.Delivery)
		}
		messageMap[date][msg.Type] = d

		if len(messageMap[date]) == 2 {
			stepsMsg := messageMap[date]["Steps"]
			activeCaloriesMsg := messageMap[date]["Active Calories"]
			// acknowledge both messages
			stepsMsg.Ack(false)
			activeCaloriesMsg.Ack(false)

			var steps IncomingMessage
			json.Unmarshal(stepsMsg.Body, &steps)
			var activeCalories IncomingMessage
			json.Unmarshal(activeCaloriesMsg.Body, &activeCalories)

			combined := OutgoingMessage{
				Steps:          steps.Value,
				ActiveCalories: activeCalories.Value,
				Date:           date,
			}

			combinedJson, err := json.Marshal(combined)
			failOnError(err, "Failed to marshal JSON")

			err = ch.Publish(
				"",             // exchange
				outgoingQ.Name, // routing key
				false,          // mandatory
				false,          // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        combinedJson,
				})
			failOnError(err, "Failed to publish a message")

			delete(messageMap, date)
		}
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
