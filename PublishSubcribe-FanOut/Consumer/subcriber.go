package main

import (
	"log"

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

	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Panicf("%s", err)
	}

	q, err := ch.QueueDeclare("", true, false, false, false, nil)
	if err != nil {
		log.Panicf("%s", err)
	}

	err = ch.QueueBind(q.Name, "", "logs", false, nil)
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
			log.Printf("%s", string(m.Body))
			m.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-blocked
}
