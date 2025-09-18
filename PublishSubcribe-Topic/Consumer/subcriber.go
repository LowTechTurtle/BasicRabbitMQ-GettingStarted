package main

import (
	"log"
	"os"

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

	err = ch.ExchangeDeclare("logs_topic", "topic", true, false, false, false, nil)
	if err != nil {
		log.Panicf("%s", err)
	}

	q, err := ch.QueueDeclare("", true, false, false, false, nil)
	if err != nil {
		log.Panicf("%s", err)
	}

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [binding_key]...", os.Args[0])
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "logs_topic", s)
		err = ch.QueueBind(
			q.Name,       // queue name
			s,            // routing key
			"logs_topic", // exchange
			false,
			nil)
		if err != nil {
			log.Panicf("%s", err)
		}
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Panicf("%s", err)
	}

	var blocked chan struct{}

	go func() {
		for m := range msgs {
			log.Printf("%s", string(m.Body))
			m.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-blocked
}
