package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func send(ch *amqp.Channel, key string, body string) {
	ch.Publish(
		"logs-ingest",
		key,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s using key %s", body, key)
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open the channel")
	defer ch.Close()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// declare exchange and queues
	err = ch.ExchangeDeclare(
		"logs-ingest", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // noWait
		amqp.Table{
			"alternate-exchange": "unrouted",
		}, // arguments
	)

	ticker := time.NewTicker(time.Second)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at ", t.UTC())
			for i := 0; i <= 4; i++ {
				key := fmt.Sprintf("syslogs-%d.foo.bar", i)
				go send(ch, key, "test")
			}
		}
	}()

	_, err = ch.QueueDeclare("syslogs", true, false, false, false, nil)
	failOnError(err, "Failed to declare queue")
	ch.QueueBind("syslogs", "syslogs-0.foo.bar", "logs-ingest", false, nil)

	NewWorker(1, ch, "syslogs")
	NewWorker(2, ch, "syslogs")
	NewWorker(3, ch, "syslogs")

	fmt.Println("Publishing messages, hit ctrl+c to exit!")
	<-sigs
}

func NewWorker(id int, ch *amqp.Channel, qName string) {
	msgs, err := ch.Consume(
		qName, // queue
		fmt.Sprintf("c-%d", id), // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	failOnError(err, "Failed to register a consumer")

	go func() {
		for msg := range msgs {
			log.Printf("[WORKER %d] Received message %s", id, msg.RoutingKey)
		}
	}()
}
