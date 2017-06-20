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

const (
	maxPriority uint8 = 10
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func increment(t amqp.Table, field string) uint8 {
	var n uint8
	n = 0
	if val, ok := t[field]; ok {
		n = val.(uint8)
	}
	n++
	t[field] = n
	return n
}

func getPriority(retries uint8) uint8 {
	var priority uint8
	if retries > maxPriority {
		priority = 0
	} else {
		priority = maxPriority - retries
	}
	return priority
}

func GenerateTTL(retries uint8) int {
	ttl := int(math.Exp2(float64(retries))) * 1000
	if ttl > 64000 {
		ttl = 64000
	}
	return ttl
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

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
		"x-max-priority":         maxPriority,
	})
	FailOnError(err, "Failed to declare queue")

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

	FailOnError(err, "Failed to register a consumer")

	go func() {
		log.Println("Consumin...")
		for msg := range msgs {
			log.Printf("[WORKER %d] Received message %s from exchange %s", os.Getpid(), msg.RoutingKey, msg.Exchange)
			if msg.Headers == nil {
				msg.Headers = amqp.Table{}
			}
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
					Priority:        msg.Priority,
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
				FailOnError(e, "Failed publishing")
			} else {
				log.Printf("[WORKER %d] Republish failed, publushing to retry queue.", os.Getpid())
				log.Printf("[WORKER %d] headers: %s.", os.Getpid(), msg.Headers)
				// republish failed, let's increment retries and publish to the retry exchange with a ttl.
				retries := increment(msg.Headers, "x-retries")
				msg.Headers["x-retries"] = increment(msg.Headers, "x-retries")

				priority := getPriority(retries)

				log.Printf("[WORKER %d] Placing in retry queue with ttl of %d and priority of %d.", os.Getpid(), ttl, priority)

				message := amqp.Publishing{
					ContentType:     msg.ContentType,
					Body:            msg.Body,
					Headers:         msg.Headers,
					ContentEncoding: msg.ContentEncoding,
					DeliveryMode:    msg.DeliveryMode,
					Priority:        priority,
					CorrelationId:   msg.CorrelationId,
					MessageId:       msg.MessageId,
					ReplyTo:         msg.ReplyTo,
					Type:            msg.Type,
					UserId:          msg.UserId,
					Expiration:      fmt.Sprintf("%d", GenerateTTL(retries)),
				}
				log.Printf("[WORKER %d] sending message %s", os.Getpid(), message)
				ch.Publish(
					"retry",
					msg.RoutingKey,
					false, // mandatory
					false, // immediate
					message,
				)
			}

		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
