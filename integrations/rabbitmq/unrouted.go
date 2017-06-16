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

func increment(t amqp.Table, field string) int64 {
	var n int64
	n = 0
	if val, ok := t[field]; ok {
		n = val.(int64)
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
		"unrouted.messages",              // queue
		fmt.Sprintf("c-%d", os.Getpid()), // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	failOnError(err, "Failed to register a consumer")

	go func() {
		log.Println("Consumin...")
		for msg := range msgs {
			log.Printf("[WORKER %d] Received message %s from exchange %s", os.Getpid(), msg.RoutingKey, msg.Exchange)
			attempts := increment(msg.Headers, "x-attempts")
			log.Printf("[WORKER %d] Attempt #%d", os.Getpid(), attempts)

			if attempts%2 != 0 {
				// publish to original exchange
				exch := msg.Exchange
				if val, ok := msg.Headers["x-original-exchange"]; ok {
					exch = val.(string)
				} else {
					msg.Headers["x-original-exchange"] = exch
				}
				log.Printf("[WORKER %d] Republishing to exchange %s.", os.Getpid(), exch)
				details := amqp.Publishing{
					ContentType:     msg.ContentType,
					Body:            msg.Body,
					Headers:         msg.Headers,
					ContentEncoding: msg.ContentEncoding,
					DeliveryMode:    msg.DeliveryMode,
					Priority:        attempts,
					CorrelationId:   msg.CorrelationId,
					MessageId:       msg.MessageId,
					ReplyTo:         msg.ReplyTo,
					Type:            msg.Type,
					UserId:          msg.UserId,
				}
				e := ch.Publish(
					exch,
					msg.RoutingKey,
					false, // mandatory
					false, // immediate
					details,
				)
				failOnError(e, "Failed publishing")
			} else {
				log.Printf("[WORKER %d] Republish failed, publushing to retry queue.", os.Getpid())
				log.Printf("[WORKER %d] headers: %s.", os.Getpid(), msg.Headers)
				// republish failed, let's increment retries and publish to the retry exchange with a ttl.
				var retries int64
				retries = 0
				if val, ok := msg.Headers["x-retries"]; ok {
					retries = val.(int64)
				}
				retries++
				msg.Headers["x-retries"] = retries

				ttl := int(math.Exp2(float64(2*retries))) * 1000
				log.Printf("[WORKER %d] Placing in retry queue with ttl of %d.", os.Getpid(), ttl)
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
						Priority:        retries,
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
