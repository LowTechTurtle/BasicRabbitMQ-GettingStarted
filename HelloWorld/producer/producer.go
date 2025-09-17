package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicf("%s", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s", err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare("hello", // name
		false, // durable
		false, // delete if unused
		false, // exclusive
		false, // no wait
		nil,   // argument
	)

	if err != nil {
		log.Panicf("%s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello World"),
	})

	if err != nil {
		log.Panicf("%s", err)
	}

	log.Print(" [x] Message Sent")
}
