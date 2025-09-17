package main

import (
	"bytes"
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

	q, err := ch.QueueDeclare("banana_queue", true, false, false, false, nil)
	if err != nil {
		log.Panicf("%s", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Panicf("%s", err)
	}

	var blocked chan struct{}

	go func() {
		for m := range msgs {
			log.Printf("Received a message: %s", m.Body)
			dotCount := bytes.Count(m.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			m.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-blocked
}
