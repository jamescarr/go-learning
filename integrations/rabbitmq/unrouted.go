package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func increment(t amqp.Table, field string) int {
	n := 0
	if val, ok := t[field].(int); ok {
		n = val
	}
	n++
	t[field] = n
	return n
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Setup our retry exchange and queue
	ch.ExchangeDeclare(
		"retry", // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // noWait
		nil,
	)

	_, err = ch.QueueDeclare("retries", true, false, false, false, amqp.Table{
		"x-dead-letter-exchange": "unrouted",
	})
	failOnError(err, "Failed to declare queue")

	ch.QueueBind("retries", "#", "retry", false, nil)

	// Consume unrouted queue
	// 1. Increment x-attempts
	// 2. if x-attempts is odd
	//   - increment x-retries
	//   - publish to the retries exch with ttl of 1000 ** x-retries
	msgs, err := ch.Consume(
		"unrouted",                       // queue
		fmt.Sprintf("c-%d", os.Getpid()), // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	failOnError(err, "Failed to register a consumer")

	go func() {
		for msg := range msgs {
			log.Printf("[WORKER %d] Received message %s from exchange %s", os.Getpid(), msg.RoutingKey, msg.Exchange)
			attempts := increment(msg.Headers, "x-attempts")

			if attempts%2 != 0 {
				// publish to original exchange
				exch := msg.Exchange
				if val, ok := msg.Headers["x-original-exchange"].(string); ok {
					exch = val
					log.Printf("[WORKER %d] Using the original exchange %s.", os.Getpid(), val)
				}
				ch.Publish(
					exch,
					msg.RoutingKey,
					false, // mandatory
					false, // immediate
					amqp.Publishing{
						ContentType:     msg.ContentType,
						Body:            msg.Body,
						Headers:         msg.Headers,
						ContentEncoding: msg.ContentEncoding,
						DeliveryMode:    msg.DeliveryMode,
						Priority:        msg.Priority,
						CorrelationId:   msg.CorrelationId,
						MessageId:       msg.MessageId,
						ReplyTo:         msg.ReplyTo,
						Type:            msg.Type,
						UserId:          msg.UserId,
					})
			} else {
				// republish failed, let's increment retries and publish to the retry exchange with a ttl.
				msg.Headers["x-original-exchange"] = msg.Exchange
				retries := 0
				if val, ok := msg.Headers["x-retries"]; ok {
					retries = val.(int)
				}
				retries++
				msg.Headers["x-retries"] = retries

				ttl := int(math.Exp2(float64(2*retries))) * 1000
				log.Println("[WORKER %d] Placing in retry queue with ttl of %d.", os.Getpid(), ttl)
				ch.Publish(
					"retry",
					msg.RoutingKey,
					false, // mandatory
					false, // immediate
					amqp.Publishing{
						ContentType:     msg.ContentType,
						Body:            msg.Body,
						Headers:         msg.Headers,
						ContentEncoding: msg.ContentEncoding,
						DeliveryMode:    msg.DeliveryMode,
						Priority:        msg.Priority,
						CorrelationId:   msg.CorrelationId,
						MessageId:       msg.MessageId,
						ReplyTo:         msg.ReplyTo,
						Type:            msg.Type,
						UserId:          msg.UserId,
						Expiration:      fmt.Sprintf("%d", ttl),
					})

			}

		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
