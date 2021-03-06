package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var (
	amqpURI = flag.String("amqp", "amqp://guest:guest@172.17.0.4:5672/", "AMQP URI")
)

var conn *amqp.Connection
var ch *amqp.Channel
var q *amqp.Queue

func AddOrderToRabbitMQ(o order) {

	var err error

	conn, err = amqp.Dial(*amqpURI)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer conn.Close()

	q, err := ch.QueueDeclare(
		"order-queue", // name
		true,          // durable ??
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare the queue")

	payload, err := json.Marshal(o)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         payload,
			Timestamp:    time.Now(),
		})
	log.Printf("Sent order %s to queue: %s", o.ID.Hex(), "order_queue")
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func initializeAmqp() {
	flag.Parse()
}
