package main

import (
	"log"

	ampq "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := ampq.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicf("%s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		log.Panicf("%s", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Panicf("%s", err)
	}

	var block chan struct{}

	go func() {
		for m := range msgs {
			log.Printf("Got a message: %s", m.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-block
}
